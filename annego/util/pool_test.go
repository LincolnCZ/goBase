package util

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestPoolConn struct {
	index int
	close int
}

func (c *TestPoolConn) Close() error {
	c.close = 1
	return nil
}

func TestPoolReuse(t *testing.T) {
	connIndex := 0
	p := &Pool{
		Dial: func() (io.Closer, error) {
			conn := &TestPoolConn{
				index: connIndex,
				close: 0,
			}
			connIndex++
			return conn, nil
		},
		MaxIdle: 5,
	}

	// 测试Get时Pool为空
	c0, _ := p.Get()
	assert.Equal(t, c0.(*TestPoolConn).index, 0)
	c1, _ := p.Get()
	assert.Equal(t, c1.(*TestPoolConn).index, 1)
	p.Put(c0, false)
	p.Put(c1, false)

	// 测试Get时Pool存在缓存
	c0, _ = p.Get()
	assert.True(t, c0.(*TestPoolConn).index < 2)
	c1, _ = p.Get()
	assert.True(t, c1.(*TestPoolConn).index < 2)
	p.Put(c0, false)
	p.Put(c1, false)

	// 测试GetNew
	c0, _ = p.GetNew()
	assert.True(t, c0.(*TestPoolConn).index == 2)
	p.Put(c1, false)

	p.Close()
}

func TestPoolLimit(t *testing.T) {
	connIndex := 0
	p := &Pool{
		Dial: func() (io.Closer, error) {
			conn := &TestPoolConn{
				index: connIndex,
				close: 0,
			}
			connIndex++
			return conn, nil
		},
		MaxIdle:   5,
		MaxActive: 10,
	}

	conns := make([]*TestPoolConn, 0)
	for i := 0; i < 10; i++ {
		c, err := p.Get()
		assert.NoError(t, err)
		conns = append(conns, c.(*TestPoolConn))
	}
	// 测试连接数到上限
	_, err := p.Get()
	assert.Equal(t, err, ErrPoolExhausted)
	assert.Equal(t, p.ActiveCount(), 10)
	assert.Equal(t, p.IdleCount(), 0)

	// 测试连空闲连接到上限
	for _, c := range conns {
		p.Put(c, false)
	}
	assert.Equal(t, p.ActiveCount(), 5)
	assert.Equal(t, p.IdleCount(), 5)

	// 测试全部关闭
	p.Close()
	assert.Equal(t, p.ActiveCount(), 0)
	assert.Equal(t, p.IdleCount(), 0)
	for _, c := range conns {
		assert.True(t, c.close == 1)
	}
}

func TestPoolTimeout(t *testing.T) {
	connIndex := 0
	p := &Pool{
		Dial: func() (io.Closer, error) {
			conn := &TestPoolConn{
				index: connIndex,
				close: 0,
			}
			connIndex++
			return conn, nil
		},
		MaxIdle:     10,
		IdleTimeout: time.Second,
	}

	conns := make([]*TestPoolConn, 0)
	for i := 0; i < 10; i++ {
		c, err := p.Get()
		assert.NoError(t, err)
		conns = append(conns, c.(*TestPoolConn))
	}
	for _, c := range conns {
		p.Put(c, false)
	}
	assert.Equal(t, p.ActiveCount(), 10)
	assert.Equal(t, p.IdleCount(), 10)

	// 没有超时时获取
	c, _ := p.Get()
	assert.True(t, c.(*TestPoolConn).index < 10)
	assert.Equal(t, p.ActiveCount(), 10)
	assert.Equal(t, p.IdleCount(), 9)
	p.Put(c, false)

	// 超时后时获取
	time.Sleep(time.Second * 2)
	c, _ = p.Get()
	assert.True(t, c.(*TestPoolConn).index == 10)
	for _, c := range conns {
		assert.Equal(t, c.close, 1)
	}
	assert.Equal(t, p.ActiveCount(), 1)
	assert.Equal(t, p.IdleCount(), 0)
	p.Put(c, false)

	p.Close()
}

func TestPoolPut(t *testing.T) {
	connIndex := 0
	p := &Pool{
		Dial: func() (io.Closer, error) {
			conn := &TestPoolConn{
				index: connIndex,
				close: 0,
			}
			connIndex++
			return conn, nil
		},
		MaxIdle: 10,
	}

	conns := make([]io.Closer, 0)
	for i := 0; i < 5; i++ {
		c, err := p.Get()
		assert.NoError(t, err)
		conns = append(conns, c)
	}

	p.Put(conns[0], false)
	assert.Equal(t, p.ActiveCount(), 5)
	assert.Equal(t, p.IdleCount(), 1)

	p.Put(conns[1], true)
	assert.Equal(t, p.ActiveCount(), 4)
	assert.Equal(t, p.IdleCount(), 1)

	p.Put(conns[2], false)
	assert.Equal(t, p.ActiveCount(), 4)
	assert.Equal(t, p.IdleCount(), 2)

	p.Close()
}

func TestPoolFilterIdle(t *testing.T) {
	connIndex := 0
	p := &Pool{
		Dial: func() (io.Closer, error) {
			conn := &TestPoolConn{
				index: connIndex,
				close: 0,
			}
			connIndex++
			return conn, nil
		},
		MaxIdle: 10,
	}

	// 生成空闲连接
	conns := make([]*TestPoolConn, 0)
	for i := 0; i < 10; i++ {
		c, err := p.Get()
		assert.NoError(t, err)
		conns = append(conns, c.(*TestPoolConn))
	}
	for _, c := range conns {
		p.Put(c, false)
	}
	assert.Equal(t, p.IdleCount(), 10)

	// 删除index > 5 的连接
	p.FilterIdle(func(conn io.Closer) bool {
		c := conn.(*TestPoolConn)
		assert.True(t, c.index < 10)
		return c.index < 5
	})
	assert.Equal(t, p.IdleCount(), 5)

	// 检查所有连接
	p.FilterIdle(func(conn io.Closer) bool {
		c := conn.(*TestPoolConn)
		assert.True(t, c.index < 5)
		return true
	})
}

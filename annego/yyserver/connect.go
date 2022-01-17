package yyserver

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"goBase/annego/packet"
)

type readBuffer struct {
	buf      []byte
	start    int
	end      int
	readsize int // 单次读取大小
	maxsize  int // 最大缓冲区大小
}

func newReadBuffer() *readBuffer {
	r := readBuffer{}
	r.start = 0
	r.end = 0
	r.readsize = 4096       // 4KB
	r.maxsize = 1024 * 1024 // 1MB
	r.buf = make([]byte, r.readsize)
	return &r
}

func (b *readBuffer) SetReadsize(s int) {
	b.readsize = s
}

func (b *readBuffer) SetMaxSize(s int) {
	b.maxsize = s
}

func (b *readBuffer) Len() int {
	return b.end - b.start
}

// ReadFrom 从IO中读取数据，返回成功读取字节数和error
func (b *readBuffer) ReadIO(conn io.Reader) (int, error) {
	var err error
	if err = b.grow(); err != nil {
		return 0, err
	}
	n, err := conn.Read(b.buf[b.end : b.end+b.readsize])
	if err != nil {
		return n, err
	}
	b.end += n
	return n, nil
}

func (b *readBuffer) Seek() []byte {
	return b.buf[b.start:b.end]
}

// HasRead 设置已读取长度，必须小于等于buffer.Len()
func (b *readBuffer) HasRead(l int) {
	if l > b.Len() {
		panic(fmt.Sprintf("readBuffer HasRead %d > %d", l, b.Len()))
	}
	b.start += l
}

// 增长缓冲区长度，保证至少有readsize的可写空间
func (b *readBuffer) grow() error {
	if b.start > 0 {
		if b.end > b.start {
			copy(b.buf, b.buf[b.start:b.end])
		}
		b.end -= b.start
		b.start = 0
	}

	if b.end+b.readsize > b.maxsize {
		return fmt.Errorf("buffer reach max size: %d", b.maxsize)
	}
	if b.end+b.readsize > len(b.buf) {
		newbuf := make([]byte, len(b.buf)*2)
		copy(newbuf, b.buf[:b.end])
		b.buf = newbuf
	}
	return nil
}

// YYConnect 单个YY协议的连接，可以用来发送接收YY协议
// 所有成员函数并发安全
type YYConnect struct {
	// UserData 可以用来保存任意的用户数据
	UserData interface{}

	conn         net.Conn
	readMut      sync.Mutex
	writeMut     sync.Mutex
	reader       *readBuffer
	writer       *bufio.Writer
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewYYConnect(conn net.Conn) *YYConnect {
	return &YYConnect{
		UserData:     nil,
		conn:         conn,
		reader:       newReadBuffer(),
		writer:       bufio.NewWriter(conn),
		readTimeout:  0,
		writeTimeout: 0,
	}
}

// SetTimeout 设置读写超时时间，只能在创建后设置一次
func (c *YYConnect) SetTimeout(readTimeout, writeTimeout time.Duration) {
	if c.readTimeout != 0 || c.writeTimeout != 0 {
		panic("YYConnect: SetTimeout again")
	}
	c.readTimeout = readTimeout
	c.writeTimeout = writeTimeout
}

func (c *YYConnect) recvFromReader(register *packet.YYRegister) (packet.Marshallable, error) {
	msg, readsize, err := register.UnmarshalBytes(c.reader.Seek())
	if err != nil {
		return nil, err
	}
	c.reader.HasRead(readsize)
	return msg, nil
}

// Recv 接收YY协议
func (c *YYConnect) Recv(register *packet.YYRegister) (packet.Marshallable, error) {
	c.readMut.Lock()
	defer c.readMut.Unlock()

	msg, err := c.recvFromReader(register)
	if err == nil {
		return msg, nil
	} else if err != packet.ErrInputNotEnough {
		return nil, err
	}

	if c.readTimeout != 0 {
		err = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
		if err != nil {
			return nil, err
		}
	}

	for {
		if _, readerr := c.reader.ReadIO(c.conn); readerr != nil {
			return nil, readerr
		}

		msg, err = c.recvFromReader(register)
		if err == packet.ErrInputNotEnough {
			continue
		}
		break
	}
	return msg, err
}

// Send 发送YY协议
func (c *YYConnect) Send(msg packet.Marshallable) error {
	c.writeMut.Lock()
	defer c.writeMut.Unlock()

	if c.writeTimeout != 0 {
		err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
		if err != nil {
			return err
		}
	}

	pack := packet.GetMarshalPack(msg)
	_, err := c.writer.Write(pack.Bytes())
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *YYConnect) Close() error {
	return c.conn.Close()
}

func (c *YYConnect) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *YYConnect) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func Dial(network, address string) (*YYConnect, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewYYConnect(c), nil
}

package util

import (
	"container/list"
	"errors"
	"io"
	"sync"
	"time"
)

// ErrPoolExhausted is returned from a pool connection reach MaxActive
var ErrPoolExhausted = errors.New("util.Pool: connection pool exhausted")

type poolConn struct {
	c io.Closer
	t time.Time
}

// Pool maintains a pool of connections.
// Code depend on redis.Pool in github.com/gomodule/redigo
type Pool struct {
	// Dial is an application supplied function for creating and configuring a
	// connection.
	//
	// The connection returned from Dial must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	Dial func() (io.Closer, error)

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	mu     sync.Mutex // mu protects the following fields
	active int        // the number of open connections in the pool
	idle   list.List  // idle connections
}

// Get connect from pool, create if empty
func (p *Pool) Get() (io.Closer, error) {
	now := time.Now()
	p.mu.Lock()

	// Prune stale connections at the back of the idle list.
	if p.IdleTimeout > 0 {
		for {
			idle := p.idle.Back()
			if idle == nil {
				break
			}
			pc := idle.Value.(*poolConn)
			if pc.t.Add(p.IdleTimeout).After(now) {
				break
			}

			p.idle.Remove(idle)
			p.mu.Unlock()
			pc.c.Close()
			p.mu.Lock()
			p.active--
		}
	}

	// Get idle connection from the front of idle list.
	for p.idle.Front() != nil {
		pc := p.idle.Front().Value.(*poolConn)
		p.idle.Remove(p.idle.Front())
		p.mu.Unlock()
		return pc.c, nil
	}

	// Handle limit
	if p.MaxActive > 0 && p.active >= p.MaxActive {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}
	p.active++
	p.mu.Unlock()

	c, err := p.Dial()
	if err != nil {
		p.mu.Lock()
		p.active--
		p.mu.Unlock()
		return nil, err
	}
	return c, nil
}

// Get connect from pool, always create new connect
func (p *Pool) GetNew() (io.Closer, error) {
	p.mu.Lock()
	// Handle limit
	if p.MaxActive > 0 && p.active >= p.MaxActive {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}
	p.active++
	p.mu.Unlock()

	c, err := p.Dial()
	if err != nil {
		p.mu.Lock()
		p.active--
		p.mu.Unlock()
		return nil, err
	}
	return c, nil
}

// Put conn into Pool
func (p *Pool) Put(conn io.Closer, forceClose bool) {
	p.mu.Lock()
	pc := &poolConn{
		c: conn,
		t: time.Now(),
	}
	if !forceClose {
		p.idle.PushFront(pc)
		if p.idle.Len() > p.MaxIdle {
			pc = p.idle.Back().Value.(*poolConn)
			p.idle.Remove(p.idle.Back())
		} else {
			pc = nil
		}
	}

	if pc != nil {
		p.mu.Unlock()
		pc.c.Close()
		p.mu.Lock()
		p.active--
	}
	p.mu.Unlock()
}

// Close releases the resources used by the pool.
func (p *Pool) Close() {
	p.mu.Lock()
	p.active = 0
	idle := list.New()
	idle.PushBackList(&p.idle)
	p.idle.Init()
	p.mu.Unlock()

	for idle.Len() > 0 {
		b := idle.Back()
		idle.Remove(b)
		pc := b.Value.(*poolConn)
		pc.c.Close()
	}
}

// ActiveCount returns the number of connections in the pool. The count
// includes idle connections and connections in use.
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

// IdleCount returns the number of idle connections in the pool.
func (p *Pool) IdleCount() int {
	p.mu.Lock()
	idle := p.idle.Len()
	p.mu.Unlock()
	return idle
}

// FilterIdle check all connect
// if check return false, close the connect
func (p *Pool) FilterIdle(check func(conn io.Closer) bool) int {
	closeConn := make([]*poolConn, 0)

	p.mu.Lock()
	ele := p.idle.Front()
	for ele != nil {
		next := ele.Next()
		pc := ele.Value.(*poolConn)
		if !check(pc.c) {
			p.idle.Remove(ele)
			closeConn = append(closeConn, pc)
			p.active--
		}
		ele = next
	}
	p.mu.Unlock()

	for _, conn := range closeConn {
		conn.c.Close()
	}
	return len(closeConn)
}

package s2s

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"sync/atomic"

	"goBase/annego/logger"
	"goBase/annego/util"
)

var ErrEmptyProxy = errors.New("s2s.S2SPool: filter proxy empty")

type FilterHandle interface {
	Filter(map[int64]*ProxyInfo) (string, int64)
}

type S2SPoolItem struct {
	ServerID int64
	Addr     string
	Conn     io.Closer
}

func (s *S2SPoolItem) Close() error {
	return s.Conn.Close()
}

type S2SPool struct {
	S2SName string

	// Dail 建立连接对象到指定地址
	Dial func(string) (io.Closer, error)

	//Filter 选择连接建立的地址
	Filter FilterHandle

	// Pool 使用的连接池，用户可修改相应配置
	// 注意：Pool.Dial 会使用上面的Dial函数，对其修改无效
	Pool *util.Pool

	mut    sync.Mutex
	proxys map[int64]*ProxyInfo
}

func NewS2SPool(name string) *S2SPool {
	p := &S2SPool{
		S2SName: name,
		Dial:    DefaultDial,
		Filter:  new(DefaultFilter),
		Pool:    new(util.Pool),
	}
	p.Pool.MaxIdle = 4
	return p
}

func (p *S2SPool) Start() error {
	if p.proxys != nil {
		return fmt.Errorf("S2SPool %s has been started", p.S2SName)
	}
	p.Pool.Dial = p.doDial
	p.proxys = make(map[int64]*ProxyInfo)

	s2sCh := Subscribe(p.S2SName)
	if s2sCh == nil {
		return fmt.Errorf("S2SPool %s subscribe fail", p.S2SName)
	}
	go p.s2sLoop(s2sCh)
	return nil
}

func (p *S2SPool) Get() (*S2SPoolItem, error) {
	item, err := p.Pool.Get()
	if err != nil {
		return nil, err
	}
	return item.(*S2SPoolItem), nil
}

func (p *S2SPool) GetNew() (*S2SPoolItem, error) {
	item, err := p.Pool.GetNew()
	if err != nil {
		return nil, err
	}
	return item.(*S2SPoolItem), nil
}

func (p *S2SPool) Put(conn *S2SPoolItem, forceClose bool) {
	p.Pool.Put(conn, forceClose)
}

func (p *S2SPool) doDial() (io.Closer, error) {
	p.mut.Lock()
	item := new(S2SPoolItem)
	item.Addr, item.ServerID = p.Filter.Filter(p.proxys)
	p.mut.Unlock()

	if item.ServerID == 0 {
		return nil, ErrEmptyProxy
	}
	conn, err := p.Dial(item.Addr)
	if err != nil {
		return nil, err
	}
	item.Conn = conn
	return item, nil
}

func (p *S2SPool) s2sLoop(s2sCh <-chan *ProxyInfo) {
	for proxyInfo := range s2sCh {
		serverID := proxyInfo.ServerID

		p.mut.Lock()
		if proxyInfo.Statu == S2SMETA_OK {
			logger.Info("[s2sLoop] add s2s proxyInfo, name %s serverID %d groupID %d statu %d",
				proxyInfo.Name, serverID, proxyInfo.GroupID, proxyInfo.Statu)
			p.proxys[serverID] = proxyInfo
		} else {
			if _, exist := p.proxys[serverID]; exist {
				logger.Info("[s2sLoop] delete s2s proxyInfo, name %s serverID %d",
					proxyInfo.Name, serverID)
				delete(p.proxys, serverID)
				p.deletePool(serverID)
			}
		}
		p.mut.Unlock()
	}
}

func (p *S2SPool) deletePool(serverID int64) {
	p.Pool.FilterIdle(func(conn io.Closer) bool {
		item := conn.(*S2SPoolItem)
		return item.ServerID != serverID
	})
}

// DefaultDial 默认建立TCP连接
func DefaultDial(addr string) (io.Closer, error) {
	return net.Dial("tcp", addr)
}

// DefaultFilter 默认地址选择方法，轮询选择所有IP地址
type DefaultFilter struct {
	index int32
}

func (f *DefaultFilter) Filter(proxy map[int64]*ProxyInfo) (addr string, serverID int64) {
	if len(proxy) == 0 {
		serverID = 0
		return
	}

	serverList := make([]int64, 0, len(proxy))
	for id := range proxy {
		serverList = append(serverList, id)
	}
	sort.Sort(util.Int64Slice(serverList))

	oldIndex := atomic.LoadInt32(&f.index)
	newIndex := oldIndex + 1
	if int(newIndex) >= len(serverList) {
		newIndex = 0
	}
	atomic.CompareAndSwapInt32(&f.index, oldIndex, newIndex)

	p, _ := proxy[serverList[newIndex]]
	addr = fmt.Sprintf("%v:%d", p.IP(), p.Port)
	serverID = p.ServerID
	return
}

package s2s

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"syscall"
	"unsafe"

	"goBase/annego/logger"
)

var gFd = -1
var gSubscribeMeta = make(map[string]chan *S2SMeta)
var gSubscribe = make(map[string]chan *ProxyInfo)
var gMyMeta unsafe.Pointer //*S2SMeta

// Init 初始化s2s
func Init(s2sname, s2skey string) bool {
	if gFd >= 0 {
		panic("s2s has been inited")
	}

	fd := initialize(s2sname, s2skey)
	if fd == -1 {
		logger.Warning("init s2s fail fd")
		return false
	}
	gFd = int(fd)
	return true
}

// Config 配置进程自身信息
func Config(iplist map[int]net.IP, port int, property map[string]string) error {
	if gFd < 0 {
		panic("s2s must init first")
	}

	jsoninfo := encodeRegInfo(iplist, port, property)
	if len(jsoninfo) == 0 {
		return errors.New("encode bson fail")
	}
	if !setMine(jsoninfo) {
		return errors.New("s2s setMine fail")
	}
	return nil
}

// SubscribeMeta 订阅裸S2S信息，用于自定义的附加数据
func SubscribeMeta(filter *SubFilter, ch chan *S2SMeta) bool {
	if gFd < 0 {
		panic("s2s must init first")
	}
	if len(filter.Name) == 0 {
		panic("s2s subscribe name empty")
	}

	if _, ok := gSubscribeMeta[filter.Name]; ok {
		return false
	}

	f := make([]SubFilter, 1)
	f[0] = *filter
	ok := subscribe(f)
	if ok {
		gSubscribeMeta[filter.Name] = ch
		return true
	} else {
		return false
	}
}

// SubscribeWithFilter 根据详细订阅信息订阅其他s2s进程
func SubscribeWithFilter(filter *SubFilter, ch chan *ProxyInfo) bool {
	if gFd < 0 {
		panic("s2s must init first")
	}
	if len(filter.Name) == 0 {
		panic("s2s subscribe name empty")
	}

	if _, ok := gSubscribe[filter.Name]; ok {
		return false
	}

	f := make([]SubFilter, 1)
	f[0] = *filter
	ok := subscribe(f)
	if ok {
		gSubscribe[filter.Name] = ch
		return true
	} else {
		return false
	}
}

// Subscribe 订阅其他s2s进程，返回chan用来接收变化数据
func Subscribe(name string) <-chan *ProxyInfo {
	if _, ok := gSubscribe[name]; ok {
		return nil
	}

	tch := make(chan *ProxyInfo)
	filter := SubFilter{
		Name:    name,
		GroupID: 0,
		S2SType: S2S_ANY_TYPE,
	}
	if !SubscribeWithFilter(&filter, tch) {
		return nil
	}
	return tch
}

// Subscribe2 订阅其他s2s进程，可以同时传入多个名称
func Subscribe2(names []string) <-chan *ProxyInfo {
	if len(names) == 0 {
		return nil
	}

	tch := make(chan *ProxyInfo)
	for _, name := range names {
		filter := SubFilter{
			Name:    name,
			GroupID: 0,
			S2SType: S2S_ANY_TYPE,
		}
		if !SubscribeWithFilter(&filter, tch) {
			return nil
		}
	}
	return tch
}

func handleS2sStatus(statu int) {
	switch statu {
	case S2S_SESSIONON, S2S_SESSIONOFF:
		break
	case S2S_SESSIONBIND:
		if atomic.LoadPointer(&gMyMeta) == nil {
			meta := GetMine()
			if meta == nil {
				logger.Warning("s2s handleS2sStatus GetMine fail")
			} else {
				logger.Info("s2s bind finish, serverId %x groupId %d", meta.ServerID, meta.GroupID)
			}
		}
	case S2S_AUTHFAILURE, S2S_ERROR:
		panic(fmt.Sprintf("s2s handleS2sStatus fatal: %d", statu))
	default:
		logger.Warning("s2s handleS2sStatus warning: %d", statu)
	}
}

func FD_SET(p *syscall.FdSet, i int) {
	p.Bits[i/64] |= 1 << (uint(i) % 64)
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}

// Start 启动S2S订阅监测
func Start() {
	if gFd < 0 {
		panic("s2s must init first")
	}

	go func() {
		for {
			timeout := syscall.Timeval{Sec: 10, Usec: 0}
			fdSet := &syscall.FdSet{}
			FD_ZERO(fdSet)
			FD_SET(fdSet, gFd)
			if _, err := syscall.Select(gFd+1, fdSet, nil, nil, &timeout); err != nil {
				logger.Info("s2s select error %v", err)
				continue
			}

			notify := pollNotify()
			handleS2sStatus(int(notify.Status))
			logger.Debug("s2s handle statu %d size %d", notify.Status, len(notify.Metas))
			for _, s2smeta := range notify.Metas {
				if ch, ok := gSubscribeMeta[s2smeta.Name]; ok {
					chout := s2smeta
					ch <- &chout
					continue
				}

				info, e := getProxyInfo(&s2smeta)
				if e != nil {
					logger.Error("s2s format meta error: %v", e)
					continue
				}

				if ch, ok := gSubscribe[info.Name]; ok {
					ch <- info
				} else {
					logger.Error("s2s unknown name: %s", info.Name)
				}
			}
		}
	}()
}

// GetMine 立刻获取一次自身信息
func GetMine() *S2SMeta {
	if gFd < 0 {
		panic("s2s must init first")
	}

	myMeta := getMine()
	if myMeta == nil {
		return nil
	}
	atomic.StorePointer(&gMyMeta, unsafe.Pointer(myMeta))
	return myMeta
}

func GetMyServerID() int64 {
	meta := (*S2SMeta)(atomic.LoadPointer(&gMyMeta))
	if meta == nil {
		return 0
	} else {
		return meta.ServerID
	}
}

func GetMyGroupID() int32 {
	meta := (*S2SMeta)(atomic.LoadPointer(&gMyMeta))
	if meta == nil {
		return 0
	} else {
		return meta.GroupID
	}
}

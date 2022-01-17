package yyserver

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"goBase/annego/logger"
	"goBase/annego/packet"
)

// 回调函数在各自连接的goroutine中执行
// 回调函数必须做到可重入

// ConnectHandle 建立连接时调用，返回false终止连接
type ConnectHandle func(*YYConnect) bool

// MessageHandle 消息处理函数，返回false终止连接
type MessageHandle func(*YYConnect, packet.Marshallable) bool

// CloseHandle 连接关闭或异常时调用, error表明具体原因
// nil 用户通过ConnechHandle或MessageHandle主动关闭
// io.ErrClosedPipe 连接goroutine外关闭连接
// io.EOF 对端关闭连接
type CloseHandle func(*YYConnect, error)

// YYServer YY协议处理服务，对应一个监听端口
// 可以设置回调函数，对划分好的YY协议进行处理
type YYServer struct {
	listener  net.Listener
	uriHandle map[uint32]MessageHandle
	register  *packet.YYRegister

	connectHandle ConnectHandle
	closeHandle   CloseHandle
}

func NewYYServer() *YYServer {
	var server YYServer
	server.listener = nil
	server.uriHandle = map[uint32]MessageHandle{}
	server.register = packet.NewYYRegister()
	return &server
}

// RegisterConnectFunc 应该在程序启动时调用
func (self *YYServer) RegisterConnectFunc(handle ConnectHandle) {
	if self.listener != nil {
		panic("YYServer is runing")
	}
	self.connectHandle = handle
}

// RegisterCloseFunc 应该在程序启动时调用
func (self *YYServer) RegisterCloseFunc(handle CloseHandle) {
	if self.listener != nil {
		panic("YYServer is runing")
	}
	self.closeHandle = handle
}

// RegisterHandle 应该在程序启动时调用，如果已经存在引起panic
func (self *YYServer) RegisterHandle(msg packet.Marshallable, handle MessageHandle) {
	if self.listener != nil {
		panic("YYServer is runing")
	}
	if !self.register.Register(msg) {
		panic(fmt.Sprintf("YYServer uri %d has register", msg.GetURI()))
	}
	self.uriHandle[msg.GetURI()] = handle
}

// Start 之前应该完成handle设置
func (self *YYServer) Start(addr string) error {
	if self.listener != nil {
		panic("YYServer is runing")
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	self.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go self.handleConnect(conn)
			} else {
				logger.Warning("accept %v error %v", addr, err)
			}
		}
	}()
	return nil
}

// StartRange 以此探测从addr开始的，trytime个端口
func (self *YYServer) StartRange(addr string, trytime int) (err error) {
	if self.listener != nil {
		panic("YYServer is runing")
	}

	addv := strings.Split(addr, ":")
	if len(addv) < 2 {
		return fmt.Errorf("addr format error %s", addr)
	}
	ip := addv[0]
	port, err := strconv.Atoi(addv[1])
	if err != nil {
		return fmt.Errorf("addr format error %s", addr)
	}
	for i := 0; i < trytime; i++ {
		useaddr := fmt.Sprintf("%s:%d", ip, port+i)
		if self.Start(useaddr) == nil {
			return nil
		}
	}
	return fmt.Errorf("try listen %s range %d fail", addr, trytime)
}

// GetListenAddr 获取监听的地址
func (self *YYServer) GetListenAddr() net.Addr {
	if self.listener == nil {
		return nil
	}
	return self.listener.Addr()
}

func (self *YYServer) handleConnect(conn net.Conn) {
	yyconn := NewYYConnect(conn)
	defer conn.Close()

	var readerr error
	if self.connectHandle != nil {
		if self.connectHandle(yyconn) == false {
			goto FIN
		}
	}

	for {
		var msg packet.Marshallable
		msg, readerr = yyconn.Recv(self.register)
		if readerr != nil {
			break
		}

		// MessageHandle返回false，主动关闭连接
		handle, _ := self.uriHandle[msg.GetURI()]
		if !handle(yyconn, msg) {
			goto FIN
		}
	}

FIN:
	// 关闭调用CloseHandle
	if self.closeHandle != nil {
		self.closeHandle(yyconn, readerr)
	}
}

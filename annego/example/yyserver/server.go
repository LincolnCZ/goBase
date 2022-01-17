// 基本的YY服务器处理实例
package main

import (
	"os"
	"time"

	"goBase/annego/logger"
	"goBase/annego/packet"
	"goBase/annego/yyserver"
)

func connect(yyconn *yyserver.YYConnect) bool {
	logger.Info("connect %s", yyconn.RemoteAddr())
	return true
}

func message(yyconn *yyserver.YYConnect, recvmsg packet.Marshallable) bool {
	msg := recvmsg.(*PTest)
	var sum uint32
	for _, i := range msg.List {
		sum += i
	}
	logger.Info("int %d str %s sum %d", msg.Int, msg.Str, sum)

	res := &PTestRes{msg.Int, msg.Str, sum}
	yyconn.Send(res)
	return true
}

func close(yyconn *yyserver.YYConnect, err error) {
	logger.Info("close addr %s reason %s", yyconn.RemoteAddr(), err)
}

func main() {
	logger.InitLog()
	server := yyserver.NewYYServer()
	server.RegisterConnectFunc(connect)
	server.RegisterCloseFunc(close)
	server.RegisterHandle(&PTest{}, message)
	if err := server.Start(os.Args[1]); err != nil {
		logger.Warning("accept error %s", err)
		return
	}
	for {
		time.Sleep(10 * time.Second)
	}
}

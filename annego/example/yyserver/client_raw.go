// 原始的 yy客户端实例，基于原始的 net.Conn
package main

import (
	"net"
	"os"

	"goBase/annego/logger"
	"goBase/annego/packet"
)

func main() {
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		logger.Warning("connect error %v", err)
		return
	}

	msg := PTest{}
	msg.Int = 12345
	msg.Str = "abcde"
	msg.List = make([]uint32, 100)
	for i := uint32(0); i < 100; i++ {
		msg.List[i] = i
	}
	sender := packet.GetMarshalPack(&msg)
	if _, err := conn.Write(sender.Bytes()); err != nil {
		logger.Warning("write error %v", err)
		return
	}

	// 这里没有粘包处理和异常判断
	recvbuf := make([]byte, 1024)
	if _, err := conn.Read(recvbuf); err != nil {
		logger.Warning("read error %v", err)
		return
	}
	up := packet.NewUnpack(recvbuf)
	up.PopHeader()
	res := PTestRes{}
	res.Unmarshal(up)

	logger.Info("recv %d %s %d", res.Int, res.Str, res.Sum)
	conn.Close()
}

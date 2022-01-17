// 基本的YY客户端处理实例
package main

import (
	"fmt"
	"math/rand"
	"os"

	"goBase/annego/logger"
	"goBase/annego/packet"
	"goBase/annego/yyserver"
)

func main() {
	packet.DefaultYYRegister.Register(new(PTest))
	packet.DefaultYYRegister.Register(new(PTestRes))

	// 建立TCP连接
	conn, err := yyserver.Dial("tcp", os.Args[1])
	if err != nil {
		logger.Warning("connect %s error: %v", os.Args[1], err)
		return
	}
	defer conn.Close()

	var index uint32 = 0
	for {
		// 生成请求
		msg := new(PTest)
		msg.Int = index
		index++
		if _, err := fmt.Scanln(&msg.Str); err != nil {
			logger.Info("scan error %v", err)
			break
		}
		msg.List = make([]uint32, 100)
		for i := uint32(0); i < 100; i++ {
			msg.List[i] = uint32(rand.Intn(100))
		}

		// 发送请求
		if err := conn.Send(msg); err != nil {
			logger.Warning("send error: %v", err)
			break
		}

		// 接收请求
		respinf, err := conn.Recv(packet.DefaultYYRegister)
		if err != nil {
			logger.Warning("recv error: %v", err)
			break
		}
		resp := respinf.(*PTestRes)
		logger.Info("recv %d %s %d", resp.Int, resp.Str, resp.Sum)
	}
}

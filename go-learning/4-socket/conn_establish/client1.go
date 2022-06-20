package main

import (
	"log"
	"net"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
}

//网络不可达或对方服务未启动
//
//$go run client1.go
//2015/11/16 14:37:41 begin dial...
//2015/11/16 14:37:41 dial error: dial tcp :8888: getsockopt: connection refused

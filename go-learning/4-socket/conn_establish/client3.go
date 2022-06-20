package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.DialTimeout("tcp", "104.236.176.96:80", 2*time.Second)
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")
}

//网络延迟较大，Dial阻塞并超时
//
//$go run client3.go
//2015/11/17 09:28:34 begin dial...
//2015/11/17 09:28:36 dial error: dial tcp 104.236.176.96:80: i/o timeout

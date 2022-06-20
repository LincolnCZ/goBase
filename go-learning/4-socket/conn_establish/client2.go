package main

import (
	"log"
	"net"
	"time"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var sl []net.Conn
	for i := 1; i < 1000; i++ {
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}

	time.Sleep(time.Second * 10000)
}

//对方服务的listen backlog满，最终的结果：
//
//$go run client2.go
//2015/11/16 21:55:44 1 :connect to server ok
//2015/11/16 21:55:44 2 :connect to server ok
//2015/11/16 21:55:44 3 :connect to server ok
//... ...
//
//2015/11/16 21:55:44 126 :connect to server ok
//2015/11/16 21:55:44 127 :connect to server ok
//2015/11/16 21:55:44 128 :connect to server ok
//
//2015/11/16 21:55:52 129 :connect to server ok
//2015/11/16 21:56:03 130 :connect to server ok
//2015/11/16 21:56:14 131 :connect to server ok
//... ...
//2015/11/16 22:03:31 128 :connect to server ok
//2015/11/16 22:04:48 129: dial error: dial tcp :8888: getsockopt: operation timed out

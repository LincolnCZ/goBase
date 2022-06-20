package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	conn.Close()
	log.Println("close ok")

	var buf = make([]byte, 32)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("read error:", err)
	} else {
		log.Printf("read % bytes, content is %s\n", n, string(buf[:n]))
	}

	n, err = conn.Write(buf)
	if err != nil {
		log.Println("write error:", err)
	} else {
		log.Printf("write % bytes, content is %s\n", n, string(buf[:n]))
	}

	time.Sleep(time.Second * 1000)
}

//$go run server1.go
//2015/11/17 17:00:51 accept a new connection
//2015/11/17 17:00:51 start to read from conn
//2015/11/17 17:00:51 conn read error: EOF
//2015/11/17 17:00:51 write 10 bytes, content is
//
//$go run client1.go
//2015/11/17 17:00:51 begin dial...
//2015/11/17 17:00:51 close ok
//2015/11/17 17:00:51 read error: read tcp 127.0.0.1:64195->127.0.0.1:8888: use of closed network connection
//2015/11/17 17:00:51 write error: write tcp 127.0.0.1:64195->127.0.0.1:8888: use of closed network connection

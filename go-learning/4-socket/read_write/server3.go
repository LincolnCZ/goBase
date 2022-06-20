//server.go

package main

import (
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		time.Sleep(10 * time.Second)
		var buf = make([]byte, 10)
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		log.Println("accept a new connection")
		go handleConn(c)
	}
}

//Socket关闭
//
//$go run client3.go hello
//2015/11/17 13:50:57 begin dial...
//2015/11/17 13:50:57 dial ok
//
//$go run server3.go
//2015/11/17 13:50:57 accept a new connection
//2015/11/17 13:51:07 start to read from conn
//2015/11/17 13:51:07 read 5 bytes, content is hello
//2015/11/17 13:51:17 start to read from conn
//2015/11/17 13:51:17 conn read error: EOF

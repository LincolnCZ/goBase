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
		var buf = make([]byte, 65536)
		log.Println("start to read from conn")
		//c.SetReadDeadline(time.Now().Add(time.Microsecond * 10))//conn read 0 bytes,  error: read tcp 127.0.0.1:8888->127.0.0.1:60763: i/o timeout
		c.SetReadDeadline(time.Now().Add(time.Microsecond * 10))
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
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

//读取操作超时
//
//$go run server4.go
//2015/11/17 14:21:17 accept a new connection
//2015/11/17 14:21:27 start to read from conn
//2015/11/17 14:21:27 conn read 0 bytes,  error: read tcp 127.0.0.1:8888->127.0.0.1:60970: i/o timeout
//2015/11/17 14:21:37 start to read from conn
//2015/11/17 14:21:37 read 65536 bytes, content is

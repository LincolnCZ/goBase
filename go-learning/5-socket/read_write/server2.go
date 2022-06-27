//server.go

package main

import (
	"log"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
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

//Socket中有部分数据
//
//$go run client2.go hi
//2015/11/17 13:30:53 begin dial...
//2015/11/17 13:30:53 dial ok
//
//$go run server2.go
//2015/11/17 13:33:45 accept a new connection
//2015/11/17 13:33:45 start to read from conn
//2015/11/17 13:33:47 read 2 bytes, content is hi
//...

//Socket中有足够数据
//
//$go run client2.go abcdefghij12345
//2015/11/17 13:38:00 begin dial...
//2015/11/17 13:38:00 dial ok
//
//$go run server2.go
//2015/11/17 13:38:00 accept a new connection
//2015/11/17 13:38:00 start to read from conn
//2015/11/17 13:38:02 read 10 bytes, content is abcdefghij
//2015/11/17 13:38:02 start to read from conn
//2015/11/17 13:38:02 read 5 bytes, content is 12345

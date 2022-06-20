package main

import (
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		//读取连接上的数据
		var buf [1024]byte
		len, err := conn.Read(buf[:])
		fmt.Println(len, err)
		//发送数据
		_, err = conn.Write([]byte("I am server!"))
		// read from the connection
		// ... ...
		// write to the connection
		//... ...
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}
}

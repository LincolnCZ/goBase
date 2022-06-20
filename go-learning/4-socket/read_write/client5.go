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
	defer conn.Close()
	log.Println("dial ok")

	data := make([]byte, 65536)
	var total int
	for {
		n, err := conn.Write(data)
		if err != nil {
			total += n
			log.Printf("write %d bytes, error:%s\n", n, err)
			break
		}
		total += n
		log.Printf("write %d bytes this time, %d bytes in total\n", n, total)
	}

	log.Printf("write %d bytes in total\n", total)
	time.Sleep(time.Second * 10000)
}

//写阻塞：
//$go run server5.go
//2015/11/17 15:07:01 accept a new connection
//2015/11/17 15:07:16 start to read from conn
//2015/11/17 15:07:16 read 60000 bytes, content is
//2015/11/17 15:07:21 start to read from conn
//2015/11/17 15:07:21 read 60000 bytes, content is
//2015/11/17 15:07:26 start to read from conn
//2015/11/17 15:07:26 read 60000 bytes, content is
//....
//
//$go run client5.go
//2015/11/17 15:07:01 write 65536 bytes this time, 720896 bytes in total
//2015/11/17 15:07:06 write 65536 bytes this time, 786432 bytes in total
//2015/11/17 15:07:16 write 65536 bytes this time, 851968 bytes in total
//2015/11/17 15:07:16 write 65536 bytes this time, 917504 bytes in total
//2015/11/17 15:07:27 write 65536 bytes this time, 983040 bytes in total
//2015/11/17 15:07:27 write 65536 bytes this time, 1048576 bytes in total
//.... ...

//写入部分数据：
//...
//2015/11/17 15:19:14 write 65536 bytes this time, 655360 bytes in total
//2015/11/17 15:19:16 write 24108 bytes, error:write tcp 127.0.0.1:62245->127.0.0.1:8888: write: broken pipe
//2015/11/17 15:19:16 write 679468 bytes in total

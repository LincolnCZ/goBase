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
		conn.SetWriteDeadline(time.Now().Add(time.Microsecond * 10))
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

//$go run client6.go
//2015/11/17 15:26:34 begin dial...
//2015/11/17 15:26:34 dial ok
//2015/11/17 15:26:34 write 65536 bytes this time, 65536 bytes in total
//... ...
//2015/11/17 15:26:34 write 65536 bytes this time, 655360 bytes in total
//2015/11/17 15:26:34 write 24108 bytes, error:write tcp 127.0.0.1:62325->127.0.0.1:8888: i/o timeout
//2015/11/17 15:26:34 write 679468 bytes in total

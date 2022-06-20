package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lucasepe/codename"
	"goBase/go-learning/tcp-stream-proto/demo2/pkg/frame"
	"goBase/go-learning/tcp-stream-proto/demo2/pkg/packet"
)

func main() {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	frameCodec := frame.NewMyFrameCodec()
	var counter int

	for {
		// send submit
		counter++
		id := fmt.Sprintf("%08d", counter) // 8 byte string
		payload := codename.Generate(rng, 4)
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}
		fmt.Printf("send submit id = %s, payload=%s\n", s.ID, s.Payload)

		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}

		err = frameCodec.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}

		// handle ack
		// read from the connection
		ackFramePayLoad, err := frameCodec.Decode(conn)
		if err != nil {
			panic(err)
		}

		p, err := packet.Decode(ackFramePayLoad)

		submitAck, ok := p.(*packet.SubmitAck)
		if !ok {
			panic("not submitack")
		}
		fmt.Printf("the result of submit ack[%s] is %d\n", submitAck.ID, submitAck.Result)

		time.Sleep(5 * time.Second)
	}
}

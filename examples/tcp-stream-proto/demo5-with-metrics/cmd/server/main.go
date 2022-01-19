package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"goBase/examples/tcp-stream-proto/demo5-with-metrics/pkg/frame"
	"goBase/examples/tcp-stream-proto/demo5-with-metrics/pkg/monitor"
	"goBase/examples/tcp-stream-proto/demo5-with-metrics/pkg/packet"
)

type customCodecServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func (cs *customCodecServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("custom codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (cs *customCodecServer) React(framePayload []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// packet decode
	var p packet.Packet
	var ackFramePayload []byte
	p, err := packet.Decode(framePayload)
	if err != nil {
		fmt.Println("react: packet decode error:", err)
		action = gnet.Close // close the connection
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			action = gnet.Close // close the connection
			return
		}
		out = []byte(ackFramePayload)
		monitor.SubmitInTotal.Add(1)
	default:
		return nil, gnet.Close // close the connection
	}

	if cs.async {
		data := append([]byte{}, ackFramePayload...)
		_ = cs.workerPool.Submit(func() {
			//fmt.Println("handleConn: async write ackFramePayload")
			c.AsyncWrite(data)
		})
		return
	}
	out = ackFramePayload
	return
}

func customCodecServe(addr string, multicore, async bool, codec gnet.ICodec) {
	var err error
	codec = frame.Frame{}
	cs := &customCodecServer{addr: addr, multicore: multicore, async: async, codec: codec, workerPool: goroutine.Default()}
	err = gnet.Serve(cs, addr, gnet.WithMulticore(multicore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}

func main() {
	// expvars http server
	go func() {
		http.ListenAndServe(":8090", nil)
	}()
	var port int
	var multicore bool

	// Example command: go run server.go --port 8888 --multicore=true
	flag.IntVar(&port, "port", 8888, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()
	addr := fmt.Sprintf("tcp://:%d", port)
	customCodecServe(addr, multicore, true, nil)
}

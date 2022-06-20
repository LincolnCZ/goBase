all: server client

server: 
	go build goBase/examples/tcp-stream-proto/demo2/cmd/server
client: 
	go build goBase/examples/tcp-stream-proto/demo2/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client

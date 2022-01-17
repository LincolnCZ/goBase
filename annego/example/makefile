.PHONY : all clean
all : client_raw client server timer s2s_example

client_raw: yyserver/client_raw.go yyserver/protocol.go
	go build -o $@ $^

client: yyserver/client.go yyserver/protocol.go
	go build -o $@ $^

server: yyserver/server.go yyserver/protocol.go
	go build -o $@ $^

timer: timer.go
	go build $^

s2s_example: s2s_example.go 
	go build $^

clean:
	rm -f client_raw client server timer s2s_example

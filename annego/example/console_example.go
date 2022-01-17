package main

import (
	"time"

	"goBase/annego/logger"
	"goBase/annego/yyserver"
)

func test1(params []string) string {
	return "test1"
}

func echo(params []string) string {
	if len(params) == 2 {
		return params[1]
	}
	return ""
}

func main() {
	console := yyserver.DefaultConsole
	console.AddCommand("test1", "print test1", test1)
	console.AddCommand("echo", "echo param", echo)
	err := console.Start("127.0.0.1:6000")
	logger.Info("start console %v", err)
	for {
		time.Sleep(1)
	}
}

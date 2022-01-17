package main

import (
	"fmt"
	"goBase/annego/yyserver"
	"time"
)

func timer3(tm time.Time) {
	fmt.Println("time 3 tick:", tm)
}

func timer5(tm time.Time) {
	fmt.Println("time 5 tick:", tm)
}

func main() {
	timer := yyserver.NewTimer()
	timer.AddHandle(3, timer3)
	timer.AddHandle(5, timer5)
	timer.Start()
}

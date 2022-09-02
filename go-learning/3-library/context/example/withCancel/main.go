package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "metrics 1")
	go watch(ctx, "metrics 2")
	go watch(ctx, "metrics 3")

	time.Sleep(10 * time.Second)
	fmt.Println("begin to stop metrics")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "metrics is stopping...")
			return
		default:
			fmt.Println(name, "goroutine monitor...")
			time.Sleep(2 * time.Second)
		}
	}
}

//func Stream(ctx context.Context, out chan<- Value) error {
//    for {
//        v, err := DoSomething(ctx)
//
//        if err != nil {
//            return err
//        }
//        select {
//        case <-ctx.Done():
//            return ctx.Err()
//        case out <- v:
//        }
//    }
//}

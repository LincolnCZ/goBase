package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func gen(done chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or done is closed, then calls
	// wg.Done.
	// 为每个输入 channel 启动一个 goroutine，将输入 channel 中的数据拷贝到
	// out channel 中，直到输入 channel，即 c，或 done 关闭。
	// 接着，退出循环并执行 wg.Done()
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

const (
	sqNum = 3
)

func main() {
	defer func() {
		time.Sleep(3 * time.Second)
		fmt.Println("the number of goroutines: ", runtime.NumGoroutine())
	}()

	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline exits, as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	defer close(done)

	in := gen(done, 2, 3, 4, 5, 6)

	// Distribute the sq work across sqNum goroutines that read from in.
	sqOut := make([]<-chan int, sqNum)
	for i := 0; i < sqNum; i++ {
		sqOut[i] = sq(done, in)
	}

	// Consume the first value from output.
	out := merge(done, sqOut...)
	//for ou := range out {
	//	fmt.Println(ou)
	//}
	fmt.Println("result:", <-out) // 4 or 9 or 16 or 25 or 36

	// done will be closed by the deferred call.
}

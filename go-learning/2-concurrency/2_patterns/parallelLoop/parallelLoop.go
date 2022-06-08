package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func handleFile(filenames []string) int64 {
	costs := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for _, f := range filenames {
		wg.Add(1)
		// worker
		go func(f string) {
			defer wg.Done()
			costTime, err := work(f)
			if err != nil {
				log.Println(err)
				return
			}

			costs <- costTime
		}(f) // 使用匿名函数，避免捕获迭代变量 f
	}

	// close和wait的操作 必须和main goroutine并行执行
	go func() {
		wg.Wait()
		close(costs)
	}()

	var total int64
	// range 循环
	for cost := range costs {
		total += cost
	}
	return total
}

////错误写法1：如果我们将wg.Wait()操作放在 range循环之前的main goroutine 中。注意，我们使用的costs通道是无缓冲的，
////   因此所有的worker goroutine都会阻塞在 costs <- cost_time，产生死锁。
////错误做法1：将等待操作放在 range 循环之前的main goroutine 中，将会产生死锁
//func handleFile(filenames []string) int64 {
//	...
//	wg.Wait()
//	close(costs)
//
//	var total int64
//	for cost := range costs {
//		total += cost
//	}
//	return total
//}

////错误写法2：如果放在循环后面，它将不可达，因为没有任何东西可用来关闭通道，循环可能永不结束。
//func handleFile(filenames []string) int64 {
//	...
//	var total int64
//	for cost := range costs {
//		total += cost
//	}
//
//	wg.Wait()
//	close(costs)
//
//	return total
//}

func work(file string) (int64, error) {
	defer fmt.Printf("finish handle file: %s\n", file)
	fmt.Printf("begin handle file: %s\n", file)

	if file == "badFile" {
		return 0, errors.New("bad file")
	}

	sleep := rand.Intn(10000)
	costTime := time.Duration(sleep) * time.Millisecond
	time.Sleep(costTime)

	return int64(costTime), nil
}

func main() {
	filenames := []string{"redFile", "blueFile", "greenFile", "badFile"}
	fmt.Printf("total cost time: %d", handleFile(filenames))
}

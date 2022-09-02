package main

import (
	"fmt"
	"time"
)

//传递给函数
func sender(ch chan string) {
	//向ch这个channel放数据
	ch <- "1"
	ch <- "2"
	ch <- "3"
	ch <- "4"
	close(ch)
}

func receiver(ch chan string) {
	////从ch这个channel读数据：
	//<-ch             // 从ch中读取一个值
	//val = <-ch
	//val := <-ch      // 从ch中读取一个值并保存到val变量中

	////检测channel是否关闭：
	//val, ok = <-ch    // 从ch读取一个值，判断是否读取成功，如果成功则保存到val变量中
	//if !ok {
	//	fmt.Println("Channel was closed")
	//}

	//range 遍历
	//  使用range来迭代channel，它会返回每次迭代过程中所读取的数据，直到channel被关闭。
	//  必须注意，只要channel未关闭，range迭代channel就会一直被阻塞。
	for val := range ch {
		fmt.Println(val)
	}
}

func main() {
	//channel --- 引用类型
	//1. 像 map 一样，channel是一个使用 make 创建的数据结构的引用。当复制或者作为参数传递到一个函数时，复制的是引用，这样调用者和被调用者都
	//   引用同一份数据结构。和其他引用类型一样，通道的零值是 nil。
	//Go 内建的函数 close、cap、len 都可以操作 chan 类型：close 会把 chan 关闭掉，cap 返回 chan 的容量，len 返回 chan 中缓存的还未被取走的元素数量。

	//创建和初始化unbuffered channel
	ch := make(chan string)
	//创建和初始化buffered channel
	//ch := make(chan string, 4)

	//goroutine
	go sender(ch)   // sender goroutine
	go receiver(ch) // receiver goroutine

	time.Sleep(1e9)

	//2. channel的send、receive、close
	//buffered channel处理:
	//|        |  nil    | empty             |full               |not full && not empty   |closed                   |
	//| -----  | -----   | -----             | -----             | -----                  | -----                   |
	//|receive |  block  |  block            |read value         |read value              |返回未读的元素，读完后返回零值|
	//|send    |  block  | write value       |block              |write value             |panic                    |
	//|close   |  panic  | closed，没有未读元素|closed，保留未读的元素|closed，保留未读的元素     |panic                    |
	//
	//unbuffered channel处理：
	//  • sender端向channel中send一个数据，然后阻塞，直到receiver端将此数据receive
	//  • receiver端一直阻塞，直到sender端向channel发送了一个数据
	//  • 关闭channel后，recv操作将获取所有已经发送的值，直到通道为空；这时任何接收操作会立即完成，同时获取到一个通道元素类型对应的零值以及一个状态码false
	//    • 利用这个特性，通过关闭 channel 实现广播操作，因为在一个已关闭的 channel 接收数据会立刻返回，并且会得到一个零值。

	//3. select 多路复用
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go worker(ch1)
	go worker(ch2)
	go stopper(ch3)

	for {
		select {
		case i := <-ch1:
			fmt.Println("Worker1 job done", i)
		case j := <-ch2:
			fmt.Println("Worker2 job done", j)
		case _, ok := <-ch3:
			if ok {
				fmt.Println("Job continue")
			} else {
				fmt.Println("Kill all job")
				return
			}
		}
	}
}

func worker(ch chan int) {
	for i := 0; i < 10000000; i++ {
		ch <- i
	}
}

func stopper(ch chan int) {
	time.Sleep(time.Second)
	ch <- 0
	time.Sleep(time.Second)
	close(ch)
}

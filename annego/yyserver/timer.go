package yyserver

import "time"

// TimerHandle 定时器回调，传人当前时间
type TimerHandle func(time.Time)

// Timer 不推荐使用，请使用 time.Ticker ！！！
// 定时器以秒为最小单位,所有TimerHandle在Timer的goroutine中执行
// 如果需要更高精度或控制执行goroutine请使用标准库time.Ticker
type Timer struct {
	interval []int
	lasttime []int64
	handle   []TimerHandle
}

func NewTimer() *Timer {
	return &Timer{}
}

// AddHandle 添加定时器函数，时间间隔sec最小1s
func (self *Timer) AddHandle(sec int, handle TimerHandle) {
	if len(self.lasttime) > 0 {
		panic("add timer handle when timer runing")
	}
	if sec < 1 {
		sec = 1
	}
	self.interval = append(self.interval, sec)
	self.handle = append(self.handle, handle)
}

func gcd(a, b int) int {
	for b != 0 {
		r := b
		b = a % b
		a = r
	}
	return a
}

// Start 启动定时器，并不断循环
func (self *Timer) Start() {
	handlesize := len(self.handle)
	if handlesize == 0 {
		return
	}

	now := time.Now()
	self.lasttime = make([]int64, handlesize)
	sleeptime := self.interval[0]

	for i := 0; i < handlesize; i++ {
		self.lasttime[i] = now.Unix()
		if sleeptime > self.interval[i] {
			sleeptime = gcd(sleeptime, self.interval[i])
		} else {
			sleeptime = gcd(self.interval[i], sleeptime)
		}
	}

	go func() {
		for {
			now = time.Now()
			nowunix := now.Unix()
			for i := 0; i < handlesize; i++ {
				if self.lasttime[i]+int64(self.interval[i]) <= nowunix {
					self.handle[i](now)
					self.lasttime[i] = nowunix
				}
			}
			time.Sleep(time.Duration(sleeptime) * time.Second)
		}
	}()
}

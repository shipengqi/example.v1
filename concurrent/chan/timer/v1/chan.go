package main

import (
	"fmt"
	"time"
)

// time.NewTimer 和 time.AfterFunc 函数可以创建 time.Timer 类型 struct
// time.Timer 方法集合中有两个方法 Stop 和 Reset，分别用于停止定时器和重置定时器
// time.Timer 的一个重要字段 C。这是一个 chan time.Time 类型的带缓冲的接收 chan。
// 一旦触及到期时间，定时器就会向这个 C 通道发送一个元素值，也就是到期时间。对应的是
// NewTimer 函数传入的那么 time.Duration 类型的值，就是定时器的到期时间
func main() {
	timer := time.NewTimer(time.Second * 2)
	fmt.Printf("present time: %v\n", time.Now())
	expireTime := <-timer.C
	fmt.Printf("expire time: %v\n", expireTime)
	fmt.Printf("stop timer: %v\n", timer.Stop())
}

// Output:
//present time: 2021-01-05 17:03:12.6184089 +0800 CST m=+0.008980401
//expire time: 2021-01-05 17:03:14.6184316 +0800 CST m=+2.008984001
//stop timer: false
//stop timer: false 是因为停止定时器时，定时器已经到期了

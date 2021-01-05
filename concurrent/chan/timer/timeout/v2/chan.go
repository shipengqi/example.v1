package main

import (
	"fmt"
	"time"
)

// 定时器是可以复用的，可以避免每次都初始化一个新的定时器，减少资源浪费
// 注意：如果在定时器到期之前，调用 Stop 方法，并返回 true，那么再从定时器的 C chan 接收元素，是不会有结果的。
// 而且会导致 goroutine 永远阻塞。所以在重置定时器之前，一定不要再对 C 执行接收操作

// 如果定时器到期，但是并没有接收 C 中的元素，那么 C 会一直缓冲这个元素，即使重置也一样。
// 因此，会影响重置后的定时器再次发送元素，虽然不会阻塞，但是该元素会被丢弃。所以要复用定时器，要确保旧值已经被接收

// 如果 NewTimer 传入的参数是负数，那么定时器会立即过期，就没有意义

// AfterFunc 接收两个参数，第一个参数是到期时间，第二个参数是到期时需要执行的函数。
// 到期时不会向 C 发送通知，而是新启用一个 goroutine 执行第二个参数传入的函数。
func main() {
	intChan := make(chan int, 1)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
			intChan <- i
		}
		close(intChan)
	}()
	timeout := time.Millisecond * 500
	var timer *time.Timer
	for {
		if timer == nil {
			timer = time.NewTimer(timeout)
		} else {
			timer.Reset(timeout)
		}
		select {
		case e, ok := <-intChan:
			if !ok {
				fmt.Println("End.")
				return
			}
			fmt.Printf("Received: %v\n", e)
		case <-timer.C:
			fmt.Println("Timeout!")
		}
	}
}

// Output:
//Timeout!
//Received: 0
//Timeout!
//Timeout!
//Received: 1
//Timeout!
//Received: 2
//Timeout!
//Received: 3
//Timeout!
//Received: 4
//End.

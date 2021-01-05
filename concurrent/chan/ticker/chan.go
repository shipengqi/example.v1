package main

import (
	"fmt"
	"time"
)

// time.Ticker 和 time.Timer 包含的字段一样，但它是一个循环定时器
// C 字段一样是一个接受通知的 chan，容量是 1。

// 如果 ticker 向 C 发送通知时，旧值未被接收，那么就会取消本次发送操作，这与定时器一致

// Ticker 只有 Stop 方法。如果 ticker 被停止，就不会再向 C 发送元素。如果这时 C 中还有元素，
// 就会一直在那里，直到被接收
func main() {
	intChan := make(chan int, 1)
	syncChan := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	go func() {
		//for _ = range ticker.C { // Stop 不会关闭 chan，所以这里会一直阻塞下去，End. [sender] 不会输出
		//	select {
		//	case intChan <- 1:
		//	case intChan <- 2:
		//	case intChan <- 3:
		//	}
		//}
		// 使用下面的代码，就可以正常输出 End. [sender]
	Loop:
		for {
			<-ticker.C
			select {
			case intChan <- 1:
			case intChan <- 2:
			case intChan <- 3:
			case <-syncChan:
				// Stop turns off a ticker. After Stop, no more ticks will be sent.
				// Stop does not close the channel, to prevent a concurrent goroutine
				ticker.Stop()
				break Loop
			}
		}
		fmt.Println("End. [sender]")
	}()
	var sum int
	var times int
	for e := range intChan {
		fmt.Printf("Received: %v\n", e)
		sum += e
		times++
		fmt.Printf("Received times: %d\n", times)
		if sum > 10 {
			fmt.Printf("Got: %v\n", sum)
			break
		}
	}
	fmt.Println("End. [receiver]")
	syncChan <- struct{}{}
}

// Output:
//Received: 1
//Received times: 1
//Received: 2
//Received times: 2
//Received: 1
//Received times: 3
//Received: 2
//Received times: 4
//Received: 3
//Received times: 5
//Received: 3
//Received times: 6
//Got: 12
//End. [receiver]

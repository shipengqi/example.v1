package main

import (
	"fmt"
	"time"
)

// cap 为 0 的 chan 是无缓冲 chan
// 同步的方式传递元素

func main() {
	sendingInterval := time.Second
	receptionInterval := time.Second * 2
	intChan := make(chan int, 0) // make(chan int) 也表示无缓冲 chan
	go func() {
		var ts0, ts1 int64
		for i := 1; i <= 5; i++ {
			intChan <- i
			ts1 = time.Now().Unix()
			if ts0 == 0 {
				fmt.Println("sent:", i)
			} else {
				fmt.Printf("sent: %d [interval: %d s]\n", i, ts1-ts0)
			}
			ts0 = time.Now().Unix()
			time.Sleep(sendingInterval)
		}
		close(intChan)
	}()
	var ts0, ts1 int64
Loop:
	for {
		select {
		case v, ok := <-intChan:
			if !ok {
				break Loop
			}
			ts1 = time.Now().Unix()
			if ts0 == 0 {
				fmt.Println("received:", v)
			} else {
				fmt.Printf("received: %d [interval: %d s]\n", v, ts1-ts0)
			}
		}
		ts0 = time.Now().Unix()
		time.Sleep(receptionInterval)
	}
	fmt.Println("end.")
}

// Output:
//sent: 1
//received: 1
//sent: 2 [interval: 2 s]
//received: 2 [interval: 2 s]
//received: 3 [interval: 2 s]
//sent: 3 [interval: 2 s]
//received: 4 [interval: 2 s]
//sent: 4 [interval: 2 s]
//received: 5 [interval: 2 s]
//sent: 5 [interval: 2 s]
//end.

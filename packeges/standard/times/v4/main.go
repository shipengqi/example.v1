package main

import (
	"fmt"
	"time"
)

func main() {
	// 初始化断续器,间隔1s
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(time.Second * 5) // 阻塞 main
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

// Output:
// Tick at 2022-01-29 10:05:14.2236825 +0800 CST m=+1.012419801
// Tick at 2022-01-29 10:05:15.2262007 +0800 CST m=+2.014937801
// Tick at 2022-01-29 10:05:16.2275316 +0800 CST m=+3.016268501
// Tick at 2022-01-29 10:05:17.22814 +0800 CST m=+4.016876701
// Ticker stopped

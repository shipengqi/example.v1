package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// 解决 time.After 内存泄露的问题
func main() {
	ch := make(chan int, 10)
	duration := 3 * time.Minute
	timer := time.NewTimer(duration)
	defer timer.Stop()

	go func() {
		var i = 1
		for {
			i++
			ch <- i
		}
	}()

	go func() {
		// 开启pprof，监听请求
		ip := "127.0.0.1:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()

	var count int
	for {
		timer.Reset(duration)
		select {
		case x := <-ch:
			fmt.Println("ch: ", x)
		case <-timer.C:
			fmt.Println("timeout: ", time.Now().Unix())
		}
		count++
		fmt.Println("count: ", count)
	}
}

// Golang中 time 包有两个定时器，分别为 ticker 和 timer。两者都可以实现定时功能，但各自都有自己的使用场景。
//
// ticker 定时器表示每隔一段时间就执行一次，一般可执行多次。
// timer 定时器表示在一段时间后执行，默认情况下只执行一次，如果想再次执行的话，每次都需要调用 time.Reset() 方法，此时效果类似 ticker 定时器。
// 同时也可以调用 stop() 方法取消定时器
// timer 定时器比 ticker 定时器多一个 Reset() 方法，两者都有 Stop() 方法，表示停止定时器,底层都调用了 stopTimer() 函数。

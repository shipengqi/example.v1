package main

import (
	"time"
)

func main()  {
	ch := make(chan int)
	go func() {
		for range time.Tick(1 * time.Second) {
			ch <- 0
		}
	}()

	for {
		select {
		case <-ch:
			println("case1")
		case <-ch:
			println("case2")
		}
	}
}

// Output:
//case1
//case1
//case2
//case1
//case2
//case2
// ...
// 上面的输出可以看出，select 在遇到多个 <-ch 同时满足可读或者可写条件时会随机选择一个 case 执行其中的代码。

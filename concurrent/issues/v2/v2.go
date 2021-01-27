package main

import (
	"fmt"
	"sync"
)

func main()  {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()
}

// 输出可能是下面这样：
//4
//2
//3
//0
//1
// 每次运行的输出都是无序的
// 原因在于，即使所有的 goroutine 都创建完了，但 goroutine 不一定已经开始运行了
// goroutine 的调度存在一定的随机性，那么其输出的结果就势必是无序且不稳定的
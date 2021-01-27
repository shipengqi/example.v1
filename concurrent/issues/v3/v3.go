package main

import (
	"fmt"
	"sync"
)

func main()  {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}

// 输出可能是下面这样：
//5
//5
//5
//2
//5
// 也可能是：
//5
//5
//5
//5
//5
// 为什么输出结果是 5 不是 4？而且并不是 100% 都是 5。
// 创建 goroutine 与真正执行 fmt.Println 并不同步。因此很有可能在你执行 fmt.Println 时，循环 for-loop 已经运行完毕，
// 因此变量 i 的值最终变成了 5。那么相反，其也有可能没运行完，存在随机性。
// goroutine 的调度存在一定的随机性，那么其输出的结果就势必是无序且不稳定的
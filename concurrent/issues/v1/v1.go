package main

import (
	"fmt"
	"sync"
)

func main()  {
	count := 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				count++
			}
		}()
	}
	wg.Wait()

	fmt.Println(count)
}

// 每次输出的都不一样，并且不会达到一百万
// 原因在于 count++ 并不是一个原子操作，在汇编上就包含了好几个动作，如下：
// MOVQ "".count(SB), AX
// LEAQ 1(AX), CX
// MOVQ CX, "".count(SB)
// 因此在执行的时候可能执行了一半就被调度系统打断，去执行别的代码。可能会同时存在多个 goroutine 同时读取到 count 的值为 1212，并各自自增 1，
// 再将其写回。
// 与此同时也会有其他的 goroutine 可能也在其自增时读到了值，形成了互相覆盖的情况，这是一种并发访问共享数据的错误。


// 使用 go run -race main.go 来做竞争检测
// 编译器会通过探测所有的内存访问，监听其内存地址的访问（读或写）。在应用运行时就能够发现对共享变量的访问和操作，进而发现问题并打印出相关的警告信息。
// 需要注意的一点是，go run -race 是运行时检测，并不是编译时。且 race 存在明确的性能开销，通常是正常程序的十倍，因此不要在生产环境打开这个配置。
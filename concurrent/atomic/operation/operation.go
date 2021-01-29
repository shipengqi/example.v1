package main

// 互斥锁是一种数据结构，用来让一个线程执行程序的关键部分，完成互斥的 "多个" 操作。
// 原子操作是针对某个值的 "单个" 互斥操作。

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var n int32
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			// 输出并不是 1000 ！！因为这个操作被编译为汇编代码后不止一条指令，
			// 因此在执行的时候可能执行了一半就被调度系统打断，去执行别的代码。
			// n ++
			// 所以要使用原子操作
			atomic.AddInt32(&n, 1)
			wg.Done()
		}()
	}
	wg.Wait()

	// fmt.Println("n: ", n)
	fmt.Println("n: ", atomic.LoadInt32(&n)) // 输出： 1000

	var n2 int32
	var lock sync.Mutex
	var wg2 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			lock.Lock()
			atomic.AddInt32(&n2, 1)
			lock.Unlock()
			wg2.Done()
		}()
	}
	wg2.Wait()
	fmt.Println("n2: ", n2) // 输出： 1000
}

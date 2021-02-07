package main

import (
	"context"
	"fmt"
	"sync/atomic"
)

// Context 类型是一种非常通用的同步工具。它的值不但可以被任意地扩散，而且还可以被用来传递额外的信息和信号。
// 更具体地说，Context 类型可以提供一类代表上下文的值。此类值是并发安全的，也就是说它可以被传播给多个 goroutine。
// Context 类型实际上是一个接口类型，而 context 包中实现该接口的所有私有类型，都是基于某个数据类型的指针类型，所以，如此传播
// 并不会影响该类型值的功能和安全。


func main() {
	coordinateWithContext()
}

func coordinateWithContext() {
	total := 12
	var num int32
	fmt.Printf("The number: %d [with context.Context]\n", num)
	ctx, cancelFunc := context.WithCancel(context.Background())
	for i := 1; i <= total; i++ {
		go addNum(&num, i, func() {
			// 原子地读取 num 变量的值，并判断它是否等于 total 变量的值
			// 如果两个值相等，那么就调用 cancelFunc 函数。也就是如果所有的 addNum 函数都执行完毕，
			// 那么就立即通知子 goroutine
			if atomic.LoadInt32(&num) == int32(total) {
				cancelFunc()
			}
		})
	}

	// Done 方法会返回一个元素类型为 struct{} 的接收通道。这个接收通道的用途是让调用方去感知“撤销”当前 Context 值的那个信号。
	// 一旦当前的 Context 值被撤销，这里的接收通道就会被立即关闭。对于一个未包含任何元素值的通道来说，它的关闭会使任何针对它的接收操作立即结束。
	<-ctx.Done()
	// Context 类型的 Err 方法的作用。该方法的结果是 error 类型的，并且其值只可能等于 context.Canceled（表示手动撤销）
	// 或者 context.DeadlineExceeded（表示由于我们给定的过期时间已到，而导致的撤销）
	// 这就可以得到“撤销”的具体原因
	fmt.Println("End.", num, ctx.Err())
}

func addNum(numP *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for {
		old := atomic.LoadInt32(numP)
		// time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, old, old+1) {
			fmt.Printf("The number: %d [%d]\n", old+1, id)
			break
		}
	}
}

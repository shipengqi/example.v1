package v2

import "sync"

// 标准库自带 sync.WaitGroup 的缺点是需要自己处理控制部分
// 代码里大量的 wg.Add 与 wg.Done 函数，所以一般封装成右侧的库
func process(stream chan int) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for elem := range stream {
				handle(elem)
			}
		}()
	}
	wg.Wait()
}

func handle(e int) {

}

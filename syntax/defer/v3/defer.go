package main

import (
	"fmt"
	"time"
)

func main() {
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))

	time.Sleep(time.Second)

	// 通过向 defer 关键字传入匿名函数来解决上面的问题
	//startedAt := time.Now()
	//defer func() { fmt.Println(time.Since(startedAt)) }()
	//
	//time.Sleep(time.Second)
}

// Output:
//0s
// 调用 defer 关键字会立刻拷贝函数中引用的外部参数，所以 time.Since(startedAt) 的结果不是在 main 函数退出之前计算的，
// 而是在 defer 关键字调用时计算的，最终导致上述代码输出 0s。
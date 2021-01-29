package main

import (
	"fmt"
)

func main()  {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			panic("panic again and again")
		}()
		panic("panic again")
	}()

	panic("panic once")
}

// Output:
//in main
//panic: panic once
//	panic: panic again
//	panic: panic again and again
//
//goroutine 1 [running]:
//main.main.func1.1()
//	D:/code/example.v1/syntax/panic/v3/panic.go:11 +0x40
//panic(0x4acea0, 0x4e97c0)
//	D:/Go/src/runtime/panic.go:969 +0x174
//main.main.func1()
//	D:/code/example.v1/syntax/panic/v3/panic.go:13 +0x62
//panic(0x4acea0, 0x4e97a0)
//	D:/Go/src/runtime/panic.go:969 +0x174
//main.main()
//	D:/code/example.v1/syntax/panic/v3/panic.go:16 +0xad
// 多次调用 panic 也不会影响 defer 函数的正常执行，所以使用 defer 进行收尾工作一般来说都是安全的。

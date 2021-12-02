package main

import (
	"fmt"
)

func main() {
	defer fmt.Println("in main")
	if err := recover(); err != nil {
		fmt.Println(err)
	}

	panic("unknown err")
}

// Output:
// panic: unknown err
//
// goroutine 1 [running]:
// main.main()
//	D:/code/example.v1/syntax/panic/v2/panic.go:13 +0x11e
// in main
// 程序没有正常退出
// recover 只有在发生 panic 之后调用才会生效。然而在上面的控制流中，recover 是在 panic 之前调用的，并不满足生效的条件，
// 所以我们需要在 defer 中使用 recover 关键字。

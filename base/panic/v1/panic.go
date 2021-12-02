package main

import "time"

func main() {
	defer println("in main")
	go func() {
		defer println("in goroutine")
		panic("")
	}()
	println("main start")
	time.Sleep(1 * time.Second)
	println("main end")
}

// Output:
// main start
// in goroutine
// panic:
//
// goroutine 5 [running]:
// main.main.func1()
//	D:/code/example.v1/syntax/panic/v1/panic.go:9 +0x78
// created by main.main
//	D:/code/example.v1/syntax/panic/v1/panic.go:7 +0x78
//
// Process finished with exit code 2
// main 函数中的 defer 语句并没有执行，执行的只有子 Goroutine 中的 defer
// defer 关键字对应的 runtime.deferproc 会将延迟调用函数与调用方所在 Goroutine 进行关联。所以当程序发生崩溃时只
// 会调用当前 Goroutine 的延迟调用函数也是非常合理的。
// 多个 Goroutine 之间没有太多的关联，一个 Goroutine 在 panic 时也不应该执行其他 Goroutine 的延迟函数。

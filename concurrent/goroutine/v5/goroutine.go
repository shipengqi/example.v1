package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("main recover", err)
		}
	}()
	fmt.Println("main start")
	go go1()
	time.Sleep(time.Second * 2)
	fmt.Println("go1 end")
	go go2()
	time.Sleep(time.Second * 2)
	fmt.Println("go2 end")
	fmt.Println("main end")
}

func go1() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("go1 recover", err)
		}
	}()
	fmt.Println("go1 start")
	panic("go1 panic")
}

func go2() {
	fmt.Println("go2 start")
	panic("go2 panic")
}

// Output:
// main start
// go1 start
// go1 recover go1 panic
// go1 end
// go2 start
// panic: go2 panic
//
// goroutine 18 [running]:
// main.go2()
//	C:/Code/example.v1/concurrent/goroutine/v5/goroutine.go:36 +0xa5
// created by main.main
//	C:/Code/example.v1/concurrent/goroutine/v5/goroutine.go:18 +0x14e
//
// Process finished with exit code 2

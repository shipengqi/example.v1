package main

import (
	"fmt"
	"sync"
)

func main()  {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i)
			wg.Done()
		}()
		// 创建该变量的副本来修复数据竞争的问题
		// go func(v int) {
		// 	fmt.Println(v)
		// 	wg.Done()
		// }(i)
	}
	wg.Wait()
}

// 此函数字面中的变量 i 与该循环中使用的是同一个变量， 因此该 Go 程中对它的读取与该递增循环产生了竞争。
// 此程序通常会打印 55555，而非 01234。
// 此程序可通过创建该变量的副本来修复
// go run -race main.go
// Output:
// ==================
// WARNING: DATA RACE
// Read at 0x00c000122068 by goroutine 7:
//  main.main.func1()
//      D:/code/example.v1/advance/race/v1/main.go:13 +0x43
//
// Previous write at 0x00c000122068 by main goroutine:
//  main.main()
//      D:/code/example.v1/advance/race/v1/main.go:11 +0x103
//
// Goroutine 7 (running) created at:
//  main.main()
//      D:/code/example.v1/advance/race/v1/main.go:12 +0xdf
// ==================
// ==================
// 2
// WARNING: DATA RACE3
//
// Read at 0x00c000122068 by goroutine 8:
//  main.main.func1()
//      D:/code/example.v1/advance/race/v1/main.go:13 +0x43
//
// Previous write at 0x00c000122068 by main goroutine:
//  main.main()
//      D:/code/example.v1/advance/race/v1/main.go:11 +0x103
//
// Goroutine 8 (running) created at:
//  main.main()
//      D:/code/example.v1/advance/race/v1/main.go:12 +0xdf
// ==================
// 3
// 5
// 4
// Found 2 data race(s)
// 可以看出 13 和 11 行有数据竞争，分别是并发的读取变量和递增变量
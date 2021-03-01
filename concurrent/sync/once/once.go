package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	var once sync.Once
	for i := 0; i < 1000; i ++ {
		once.Do(func() {
			count ++
		})
	}

	fmt.Println("count: ", count) // 1
}

// Output:
// count:  1
// 说明制运行了一次 once.Do 的函数参数
// func() {
//     count ++
// }
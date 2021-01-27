package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func() {
			fmt.Printf("Hello, %s!\n", name)
		}()
	}
	time.Sleep(time.Millisecond)
}

// 大概率会输出下面这样
// Hello, Mark!
// Hello, Mark!
// Hello, Mark!
// Hello, Mark!
// Hello, Mark!
// 因为 go 函数是在 for 循环执行完之后执行，也就是 name 被最终赋值为 Mark。
// 这不是一定的，go 函数也可能在 for 循环执行一半的时候执行

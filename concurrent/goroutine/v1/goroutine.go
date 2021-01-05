package main

import (
	"fmt"
	"time"
)

func main() {
	name := "Eric"
	go func() {
		fmt.Printf("Hello, %s!\n", name)
	}()
	name = "Harry"
	time.Sleep(time.Millisecond)
}

// 大概率会输出 Hello, Harry，因为 name = "Harry" 执行之后 go 函数才执行

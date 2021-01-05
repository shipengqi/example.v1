package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go func(name string) {
			fmt.Printf("Hello, %s!\n", name)
		}(name)
	}
	time.Sleep(time.Millisecond*2)
}

// 输出顺序不是一定的
// Hello, Eric!
// Hello, Mark
// Hello, Robert
// Hello, Harry
// Hello, Jim
// 之所以可以输出所有元素，是因为 name 是 string 类型，在作为函数参数传递时，会被复制
// 但是对于引用类型，相当于只是复制了指针，函数中的修改会影响外部的值

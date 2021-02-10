package main

import (
	"container/list"
	"fmt"
)

// e 是一个 Element 类型的指针，当然其也可能为 nil，但是 golang 中 list 包中函数没有对其进行是否为 nil 的检查，默认其对非 nil 进行操作，所
// 以这种情况下，便可能出现 runtime panic
func main()  {
	l := list.New()
	l.PushBack(1)
	fmt.Println(l.Front().Value) // 1
	value := l.Remove(l.Front())
	fmt.Println(value)            // 1
	value1 := l.Remove(l.Front()) // panic: runtime error: invalid memory address or nil pointer dereference
	fmt.Println(value1)
}

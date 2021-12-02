package main

import (
	"container/ring"
	"fmt"
)

func main() {
	r := ring.New(10) // 初始长度10
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}
	for i := 0; i < r.Len(); i++ {
		fmt.Println("[Initialized]: ", r.Value)
		r = r.Next()
	}
	r = r.Move(6)
	fmt.Println("[Moved]: ", r.Value) // 6
	r1 := r.Unlink(19)   // 移除 19%10=9 个元素
	for i := 0; i < r1.Len(); i++ {
		fmt.Println("[Unlinked]: ", r1.Value)
		r1 = r1.Next()
	}
	fmt.Println("[r length]: ", r.Len())  // 10-9=1
	fmt.Println("[r1 length]: ", r1.Len()) // 9
}

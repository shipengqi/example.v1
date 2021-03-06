package main

import (
	"container/ring"
	"fmt"
)

func main() {
	r := ring.New(3)

	for i := 1; i <= 3; i++ {
		r.Value = i
		r = r.Next()
	}

	// 计算 1+2+3
	s := 0
	r.Do(func(p interface{}) {
		s += p.(int)
	})
	fmt.Println("sum is", s)
}

// Output:
// sum is 6

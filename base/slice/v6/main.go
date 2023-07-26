package main

import "fmt"

// 值传递

func main() {
	sl := make([]int, 0, 10)
	var appendFunc = func(s []int) {
		s = append(s, 10, 20, 30)
	}
	appendFunc(sl)
	// slice 的底层存储是一个数组，slice 本身的数据结构
	// type SliceHeader struct {
	//     Data uintptr
	//     Len  int
	//     Cap int
	// }
	// 因为 go 是值传递，所以 appendFunc 修改的是 sl 的副本，len 和 cap 并没有被修改，下面的输出是 []
	fmt.Println(sl) // []
	// Data 是指向数组的指针，所以 appendFunc 可以修改底层的数组，下面的输出会包含 10 20 30
	fmt.Println(sl[:10]) // [10 20 30 0 0 0 0 0 0 0]
	// 为什么 sl[:10] 和 sl[:] 的输出不同，是因为 go 的切片的一个优化
	// slice[low:high] 中的 high，最大的取值范围对应着切片的容量（cap），不是单纯的长度（len）。
	// sl[:10] 可以输出容量范围内的值，并且没有越界。
	// sl[:] 由于 len 为 0，并且没有指定最大索引。high 则会取 len 的值，所以输出为 []
	fmt.Println(sl[:]) // []

	sl2 := []int{1, 2, 3, 4, 5}
	sl3 := sl2[1:]
	fmt.Println(sl3)
}

// Output
// []
// [10 20 30 0 0 0 0 0 0 0]
// []
// [2 3 4 5]

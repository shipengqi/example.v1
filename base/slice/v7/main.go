package main

import "fmt"

// 值传递

func main() {
	sl := make([]int, 0, 10)
	var appendFunc = func(s []int) {
		// slice 扩容
		// 当切片的容量不足时，go 会调用 runtime.growslice 函数为切片扩容
		// 如果期望容量大于当前容量的两倍就会使用期望容量
		// 如果当前切片的长度小于 1024 就会将容量翻倍
		// 如果当前切片的长度大于 1024 就会每次增加 25% 的容量，直到新容量大于期望容量
		// 如果计算新容量时发生了内存溢出或者请求内存超过上限，就会直接崩溃退出程序
		// runtime.growslice 函数最终会返回一个新的切片
		// 其中包含了新的数组指针、大小和容量，这个返回的三元组最终会覆盖原切片。
		// 所以这里的 s 被赋值了一个完整的新的切片，和外面的 sl 已经没有任何关系。
		s = append(s, 10, 20, 30, 10, 20, 30, 10, 20, 30, 10, 20, 30)
		fmt.Println(s) // [10 20 30 10 20 30 10 20 30 10 20 30]
		fmt.Println(len(s)) // 12
		fmt.Println(cap(s)) // 20
	}
	appendFunc(sl)
	fmt.Println(sl) // []

	fmt.Println(sl[:10]) // [0 0 0 0 0 0 0 0 0 0]
	fmt.Println(sl[:]) // []
	fmt.Println(len(sl)) // 0
	fmt.Println(cap(sl)) // 10
	for k := range sl { // 由于 len 是 0，所以下面的代码不会有输出
		fmt.Println(sl[k])
	}
}

// Output
// [10 20 30 10 20 30 10 20 30 10 20 30]
// 12
// 20
// []
// [0 0 0 0 0 0 0 0 0 0]
// []
// 0
// 10

// func growslice(et *_type, old slice, cap int) slice {
//	newcap := old.cap
//	doublecap := newcap + newcap
//	if cap > doublecap {
//		newcap = cap
//	} else {
//		if old.len < 1024 {
//			newcap = doublecap
//		} else {
//			for 0 < newcap && newcap < cap {
//				newcap += newcap / 4
//			}
//			if newcap <= 0 {
//				newcap = cap
//			}
//		}
//	}

package main

import "fmt"

func main()  {
	// 在遍历数组的同时修改数组的元素，能否得到一个永远都不会停止的循环
	arr := []int{1, 2, 3}
	for _, v := range arr {
		arr = append(arr, v)
	}
	fmt.Println(arr)
}

// Output:
// [1 2 3 1 2 3]
// 上面的输出意味着循环只遍历了原始切片中的三个元素，在遍历切片时追加的元素不会增加循环的执行次数
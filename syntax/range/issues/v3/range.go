package main

func main()  {
	arr := []int{1, 2, 3}
	for i, _ := range arr {
		arr[i] = 0
	}
}

// 遍历清空数组
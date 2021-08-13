package main

import "fmt"

func main()  {
	arr := []int{1, 2, 3}
	var newArr []*int
	for _, v := range arr {
		newArr = append(newArr, &v)
	}
	for _, v := range newArr {
		fmt.Println(*v)
	}

	fmt.Println("=====================")
	var newArr2 []*int
	for i := range arr {
		newArr2 = append(newArr2, &arr[i])
	}
	for _, v := range newArr2 {
		fmt.Println(*v)
	}
}

// Output:
// 3 3 3
// 正确的做法应该是使用 &arr[i] 替代 &v
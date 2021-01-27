package main

import "fmt"

func main()  {
	hash := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	for k, v := range hash {
		fmt.Println(k, v)
	}
}

// Output:
//3 3
//1 1
//2 2
// 顺序是随机的
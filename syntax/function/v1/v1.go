package main

import "fmt"

// Golang 可变参数
func Add(values ... int)  {
	for _, v := range values {
		fmt.Println(v)
	}
}

func main()  {
	Add(1, 2, 3, 4, 6, 8)

	fmt.Println("=================")
	vs := []int{
		5,
		6,
		7,
		8,
	}
	Add(vs...)
}

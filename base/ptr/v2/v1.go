package main

import "fmt"

func swap(a, b *int) {
	fmt.Printf("a value: %d\n", *a) // a value: 1
	fmt.Printf("b value: %d\n", *b) // b value: 2

	// 取 a 指针的值, 赋给临时变量 t
	t := *a
	fmt.Printf("t value: %v\n", t) // t value: 1

	// 取 b 指针的值, 赋给 a 指针指向的变量
	*a = *b
	fmt.Printf("a new value: %d\n", *a) // a new value: 2

	// 将 a 指针的值赋给 b 指针指向的变量
	*b = t
	fmt.Printf("b new value: %d\n", *b) // b new value: 1
}

func swap2(a, b *int) {
	fmt.Printf("swap2 a value: %d\n", *a) // swap2 a value: 1
	fmt.Printf("swap2 b value: %d\n", *b) // swap2 b value: 2
	// a 和 b 都是指针的拷贝
	// 交换只是 指针的拷贝指向了新的值，不会影响到函数外面的值
	b, a = a, b

	fmt.Printf("swap2 a new value: %d\n", *a) // swap2 a new value: 2
	fmt.Printf("swap2 b new value: %d\n", *b) // swap2 b new value: 1
}

func main() {
	x, y := 1, 2
	swap(&x, &y)
	fmt.Println(x, y) // 2 1

	m, n := 1, 2
	swap2(&m, &n)
	fmt.Println(m, n) // 1 2
	// swap2 交换是不成功的
}

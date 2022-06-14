package main

import "fmt"

func main() {
	// 准备一个字符串类型
	var house = "Point 10880, 90265"
	// 对字符串取地址
	ptr := &house
	// ptr 的类型
	fmt.Printf("ptr type: %T\n", ptr) // ptr type: *string
	// ptr 的指针地址
	fmt.Printf("address: %p\n", ptr) // address: 0xc00004a230
	// 对指针进行取值操作
	value := *ptr
	// 取值后的类型
	fmt.Printf("value type: %T\n", value) // value type: string
	// 指针取值后就是指向变量的值
	fmt.Printf("value: %s\n", value) // value: Point 10880, 90265
}

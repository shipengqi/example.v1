package main

import "fmt"

// func copy(dst, src []Type) int
// copy() 可以将一个数组切片复制到另一个数组切片中，**如果加入的两个数组切片不一样大，就会按照其中较小的那个数组切片的元素个数进行复制**。
// 切片 dst 需要先初始化长度
// 源切片中元素类型为引用类型时，拷贝的是引用

func main()  {
	src := []string{
		"str1",
		"str2",
		"str3",
	}

	target1 := make([]string, 0)
	target2 := make([]string, 1)
	target3 := make([]string, len(src))

	copy(target1, src)
	copy(target2, src)
	copy(target3, src)
    fmt.Println(target1)            // []
	fmt.Println(target2)            // [str1]
	fmt.Println(target3)            // [str1 str2 str3]

	// matA 和 matB 地址不一样，但 matA[0] 和 matB[0] 的地址是一样的
	matA := [][]int{
		{0, 1, 1, 0},
		{0, 1, 1, 1},
		{1, 1, 1, 0},
	}
	matB := make([][]int, len(matA))
	copy(matB, matA)
	fmt.Printf("%p, %p\n", matA, matA[0]) // 0xc000058050, 0xc000010240
	fmt.Printf("%p, %p\n", matB, matB[0]) // 0xc0000580a0, 0xc000010240

	// 正确的拷贝一个多维数组
	// 想 copy 多维切片中的每一个切片类型的元素，那么需要将每个切片元素进行 初始化 并 拷贝。注意是两步：先 初始化，再 拷贝。
	for i := range matA {
		matB[i] = make([]int, len(matA[i])) // 注意初始化长度
		copy(matB[i], matA[i])
	}
	fmt.Printf("%p, %p\n", matA, matA[0]) // 0xc000058050, 0xc000010240
	fmt.Printf("%p, %p\n", matB, matB[0]) // 0xc0000580a0, 0xc0000102a0
}

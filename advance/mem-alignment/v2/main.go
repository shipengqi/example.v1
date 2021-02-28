package main

import (
	"fmt"
	"unsafe"
)

// Sizeof Sizeof 返回类型 x 所占据的字节数，但不包含 x 所指向的内容的大小。
// 例如，对于一个指针，函数返回的大小为 8 字节（64位机上），
// 一个 slice 的大小则为 slice header 的大小。
//
// Offsetof 返回结构体成员在内存中的位置离结构体起始处的字节数，所传参数必须是结构体的成员。
//
// Alignof 返回 m，m 是指当类型进行内存对齐时，它分配到的内存地址能整除 m。

func main()  {
	type SizeOfA struct {
		A int
	}
	fmt.Println("SizeOfA: ", unsafe.Sizeof(SizeOfA{0}))

	type SizeOfB struct {
		A byte  // 1 字节
		B int32 // 4 字节
	}

	fmt.Println("SizeOfB: ", unsafe.Sizeof(SizeOfB{0, 0}))
	fmt.Println("SizeOfB Alignof: ", unsafe.Alignof(SizeOfB{0, 0}))

	// SizeOfE 中，元素的大小分别为 1,8,1，但是实际结构体占 24 字节，远超元素实际大小，
	// 因为内存对齐原因，最开始分配的8字节中包含了 1 字节的 A，剩余的 7 字节不足以放下 B，
	// 又为 B 分配了 8 字节，剩余的 C 独占再分配的 8 字节。
	type SizeOfC struct {
		A byte  // 1 字节
		B int64 // 8 字节
		C byte  // 1 字节
	}

	fmt.Println("SizeOfC: ", unsafe.Sizeof(SizeOfC{0, 0, 0}))

	type SizeOfD struct {
		A byte
		B [5]int32
	}
	fmt.Println("SizeOfD: ", unsafe.Sizeof(SizeOfD{}))
	fmt.Println("SizeOfD Alignof: ", unsafe.Alignof(SizeOfB{0, 0}))

	type SizeOfF struct {
		A byte
		C int16
		B int64
		D int32
	}

	fmt.Println("SizeOfF.A: ", unsafe.Offsetof(SizeOfF{}.A))
	fmt.Println("SizeOfF.C: ", unsafe.Offsetof(SizeOfF{}.C))
	fmt.Println("SizeOfF.B: ", unsafe.Offsetof(SizeOfF{}.B))
	fmt.Println("SizeOfF.D: ", unsafe.Offsetof(SizeOfF{}.D))
}

// Output:
// SizeOfA:  8
// SizeOfB:  8
// SizeOfB Alignof:  4
// SizeOfC:  24
// SizeOfD:  24
// SizeOfD Alignof:  4
// SizeOfF.A:  0
// SizeOfF.C:  2
// SizeOfF.B:  8
// SizeOfF.D:  16

// SizeOfF 的内存不糊：
// |A|-|C  |-|-|-|-|
// |B              |
// |D      |-|-|-|-|
// A 占用 1 个字节，C 2 个字节，B 8 个字节，D 4 个字节
// - 代表内存空洞
// 可以看出，结构体中元素不同顺序的排列会导致内存分配的极大差异，不好的顺序会产生许多的内存空洞，造成大量内存浪费。
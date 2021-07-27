package main

import (
	"fmt"
	"unsafe"
)

type Part1 struct {
	a bool
	b int32
	c int8
	d int64
	e byte
}

type Part2 struct {
	e byte
	c int8
	a bool
	b int32
	d int64
}

func main()  {
	fmt.Printf("bool size: %d\n", unsafe.Sizeof(true))
	fmt.Printf("int32 size: %d\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int8 size: %d\n", unsafe.Sizeof(int8(0)))
	fmt.Printf("int64 size: %d\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("byte size: %d\n", unsafe.Sizeof(byte(0)))
	fmt.Printf("string size: %d\n", unsafe.Sizeof("example"))

	part1 := Part1{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))

	part2 := Part2{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part2), unsafe.Alignof(part2))
}

// Output:
// bool size: 1
// int32 size: 4
// int8 size: 1
// int64 size: 8
// byte size: 1
// string size: 16
// 看上去 Part1 这一个结构体的占用内存大小为 1+4+1+8+1 = 15 个字节，但是最后一行输出
// part1 size: 32, align: 8
// 内存对齐后，占用了 32 个字节
// Part2 调整了 Part1 的字段顺序
// part1 size: 16, align: 8
// 通过调整结构体内成员变量的字段顺序达到缩小结构体占用大小
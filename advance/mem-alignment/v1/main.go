package main

import (
	"fmt"
	"unsafe"
)

func main()  {
	var vint int = 1
	var vint8 int8 = 1
	var vint16 int16 = 1
	var vint32 int32 = 1
	var vint64 int64 = 1

	fmt.Println("int: ", unsafe.Sizeof(vint))
	fmt.Println("int8: ", unsafe.Sizeof(vint8))
	fmt.Println("int16: ", unsafe.Sizeof(vint16))
	fmt.Println("int32: ", unsafe.Sizeof(vint32))
	fmt.Println("int64: ", unsafe.Sizeof(vint64))

	// var vbool bool = true
	var vstring = "a"
	var vstring2 = "ab"
	var vstring3 = "你"
	var vstring4 = "你好"
	var vstring5 = "Golang 的内存对齐 example"
	var vbyte byte = 0

	fmt.Println("bool: ", unsafe.Sizeof(true))
	fmt.Println("byte: ", unsafe.Sizeof(vbyte))
	fmt.Println("string a: ", unsafe.Sizeof(vstring))
	fmt.Println("string ab: ", unsafe.Sizeof(vstring2))
	fmt.Println("string 你: ", unsafe.Sizeof(vstring3))
	fmt.Println("string 你好: ", unsafe.Sizeof(vstring4))
	fmt.Println("string Golang 的内存对齐 example: ", unsafe.Sizeof(vstring5))
}

// Output:
// int:  8
// int8:  1
// int16:  2
// int32:  4
// int64:  8
// bool:  1
// byte:  1
// string a:  16
// string ab:  16
// string 你:  16
// string 你好:  16
// int 占用一个机器字节
// string 占用两个机器字节
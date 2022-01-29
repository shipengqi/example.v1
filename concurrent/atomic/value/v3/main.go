package main

import (
	"fmt"
	"unsafe"
)

// atomic.Value, 此类型相当于一个容器，被用来"原子地"存储（Store）和加载任意类型的值
// // A Value provides an atomic load and store of a consistently typed value.
// // The zero value for a Value returns nil from Load.
// // Once Store has been called, a Value must not be copied.
// //
// // A Value must not be copied after first use.
// type Value struct {
//	v interface{}
// }
//
// // ifaceWords is interface{} internal representation.
// type ifaceWords struct {
//	typ  unsafe.Pointer // 原始类型
//	data unsafe.Pointer // 真正的值
// }
//
// unsafe.Pointer
// Go语言并不支持直接操作内存，但是它的标准库提供一种不保证向后兼容的指针类型unsafe.Pointer，
// 让程序可以灵活的操作内存，它的特别之处在于：可以绕过Go语言类型系统的检查
// 也就是说：如果两种类型具有相同的内存结构，我们可以将unsafe.Pointer当作桥梁，让这两种类型的指针相互转换，从而实现同一份内存拥有两种解读方式


func main() {
	var a int32 = 15
	// 获得a的*int类型指针
	ptr := (*int)(unsafe.Pointer(&a))
	fmt.Println(ptr, *ptr) // 0xc00000a0b8 15
}


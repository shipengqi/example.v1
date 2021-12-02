package main

// #include <hello.h>
import "C"
import "fmt"

// 通过 go 实现 C 函数，并导出
// //export SayHello 指令将 Go 语言实现的 SayHello 函数导出为 C 函函数
// 必须是 //export // 和 export 之间不能有空格


//export SayHello
func SayHello(str *C.char) {
	fmt.Println("Go .....")
	fmt.Println(C.GoString(str))
}
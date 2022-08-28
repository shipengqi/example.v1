package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed hello1.txt hello2.txt
var f embed.FS

// 多个文件的另一种方式
// //go:embed hello1.txt
// //go:embed hello2.txt
// var f embed.FS

func main() {
	data1, _ := f.ReadFile("hello1.txt")
	fmt.Println(string(data1))

	data2, _ := f.ReadFile("hello2.txt")
	fmt.Println(string(data2))
}

// hello1, this is embed
//
// hello2, this is embed

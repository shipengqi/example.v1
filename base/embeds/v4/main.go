package main

import (
	"embed"
	_ "embed"
	"fmt"
)

//go:embed hello
var f embed.FS //  embed.FS 实现了 io/fs 接口, embed 所打包进二进制文件的内容只允许读取，不允许变更

// 能支持贪婪模式的匹配
// //go:embed helloworld/*
// var f embed.FS

func main() {
	data1, _ := f.ReadFile("hello/hello1.txt")
	fmt.Println(string(data1))

	data2, _ := f.ReadFile("hello/hello2.txt")
	fmt.Println(string(data2))
}

// hello1, this is embed
//
// hello2, this is embed

package main

import (
	"fmt"
	"os"
)

func main() {
	reader, writer, err := os.Pipe() // 生成命名管道的 reader 和 writer
	if err != nil {
		fmt.Println("Pipe", err)
		return
	}
	_, err = writer.Write([]byte("pipe content")) // writer 写数据
	if err != nil {
		fmt.Println("Write", err)
		return
	}

	// writer 只能写，reader 只能读，反向应用将出错
	// r.Write([]byte("test")) // 0: bad file descriptor

	buf := make([]byte, 20)
	n, err := reader.Read(buf) // reader 读数据
	if err != nil {
		fmt.Println("Read", err)
		return
	}
	fmt.Printf("%q\n", string(buf[:n])) // "pipe content"
}

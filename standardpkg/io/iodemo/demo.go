package main

import (
	"fmt"
	"os"
	"strings"
)

func readerDemo()  {
	reader := strings.NewReader("example.v1")
	p := make([]byte, 6)
	n, err := reader.ReadAt(p, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d\n", p, n) // ample., 6
}

func writerDemo()  {
	file, err := os.Create("writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("hello, overwrite")
	n, err := file.WriteAt([]byte("example.v1"), 7)
	if err != nil {
		panic(err)
	}
	fmt.Println(n) // 10
}
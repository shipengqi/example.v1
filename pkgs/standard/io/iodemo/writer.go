package main

import (
	"bytes"
	"fmt"
	"os"
)

// JavaScript 中单引号和双引号可以同时使用，都可以用来表示字符串。
// Java 中单引号表示 char 类型，双引号表示 string 类型。
// Go 中，双引号是用来表示 string 类型，本质是一个 []byte 类型，单引号表示 rune 类型。反引号，用来创建原生的字符串字面量，它可以由多行组成，
// 但不支持任何转义序列。因此，当把两个不同类型的变量进行拼接时，就会报错。

func writerExample() {
LOOP:
	for {
		writerMenu()
		var selected string
		_, _ = fmt.Scanln(&selected)
		switch selected {
		case "1":
			fmt.Println("Please enter a character:")
			var ch byte
			_, _ = fmt.Scanf("%c\n", &ch)
			buffer := new(bytes.Buffer)
			err := buffer.WriteByte(ch)
			if err != nil {
				fmt.Printf("write err: %v\n", err)
				continue
			}
			fmt.Println("A byte has been written, read it")
			b, _ := buffer.ReadByte()
			fmt.Printf("read the：%c\n", b)
		case "b":
			fmt.Println("back to menu!")
			break LOOP
		case "q":
			fmt.Println("quit!")
			os.Exit(0)
		default:
			fmt.Printf("unknown option: %s\n", selected)
			continue
		}

	}
}

func writerMenu() {
	fmt.Println("")
	fmt.Println("********* byte.Writer demo *********")
	fmt.Println("******* please input your option: *********")
	fmt.Println("1, stdin")
	fmt.Println("b, back to menu")
	fmt.Println("q, exit")
	fmt.Println("***********************************************")
}

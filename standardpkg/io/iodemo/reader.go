package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func readerExample() {
LOOP:
	for {
		readerMenu()

		var selected string
		_, _ = fmt.Scanln(&selected)

		var data []byte
		var err error

		switch strings.ToLower(selected) {
		case "1":
			fmt.Println("Please enter less than 9 characters, ending with Enter:")
			data, err = ReadFrom(os.Stdin, 11)
		case "2":
			filename := "README.md"
			file, err := os.Open(filename)
			if err != nil {
				fmt.Printf("open: %s, err: %v\n", filename, err)
				continue
			}
			data, err = ReadFrom(file, 9)
			_ = file.Close()
		case "3":
			data, err = ReadFrom(strings.NewReader("from string"), 12)
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

		if err != nil {
			fmt.Printf("read err: %v\n", err)
			continue
		}
		fmt.Printf("data：%s\n", data)
	}
}

// ReadForm read data from reader, return []byte
func ReadFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

func readerMenu() {
	fmt.Println("")
	fmt.Println("*******read data from kinds of sources*********")
	fmt.Println("*******please input your option: *********")
	fmt.Println("1, stdin")
	fmt.Println("2, file")
	fmt.Println("3, string")
	fmt.Println("b, back to menu")
	fmt.Println("q, exit")
	fmt.Println("***********************************************")
}

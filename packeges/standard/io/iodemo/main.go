package main

import (
	"fmt"
	"github.com/shipengqi/example.v1/utils"
)

func main()  {
	utils.Welcome()
	MainMenu()
}

func MainMenu() {
LOOP:
	for {
		fmt.Println("")
		fmt.Println("******* please input your option: *********")
		fmt.Println("1, io.Reader demo")
		fmt.Println("2, io.ByteReader/ByteWriter demo")
		fmt.Println("q, quit")
		fmt.Println("***********************************")

		var ch string
		fmt.Scanln(&ch)

		switch ch {
		case "1":
			readerExample()
		case "2":
			writerExample()
		case "q":
			fmt.Println("quit!")
			break LOOP
		default:
			fmt.Println("unknown option!")
			continue
		}
	}
}
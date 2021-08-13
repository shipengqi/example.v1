package main

import (
	"fmt"
	"strconv"
)

func makeCakeAndSend(cs chan string, count int) {
	for i := 1; i <= count; i++ {
		cakeName := "Strawberry Cake " + strconv.Itoa(i)
		cs <- cakeName //send a strawberry cake
	}
	close(cs)
}

func receiveCakeAndPack(cs chan string) {
	for s := range cs {
		fmt.Println("Packing received cake: ", s)
		//if s == "Strawberry Cake 3" {
		//	break
		//}
	}
}

func main() {
	cs := make(chan string)
	go makeCakeAndSend(cs, 5)
	// go receiveCakeAndPack(cs)
	for s := range cs {
		fmt.Println("Packing received cake: ", s)
	}
}

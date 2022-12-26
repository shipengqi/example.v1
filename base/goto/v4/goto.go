package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second)
		c1 <- "h"

		time.Sleep(time.Second)
		c1 <- "hello"

		time.Sleep(time.Second)
		c1 <- "h2"

	}()
	for {
	SELECT:
		fmt.Println("test")
		select {
		case str := <-c1:
			if str == "hello" {
				fmt.Println("goto")
				goto SELECT
			}
			fmt.Println("c1: ", str)
		case str := <-c2:
			if str == "world" {
				goto SELECT
			}
			fmt.Println("c2: ", str)
		}
		fmt.Println("sleep")
		time.Sleep(time.Second)
	}
}

// Output：
// test
// c1:  h
// sleep
// test
// goto
// test
// c1:  h2
// sleep
// test
// fatal error: all goroutines are asleep - deadlock!

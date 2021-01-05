package main

import (
	"fmt"
	"time"
)

// time.After(time.Millisecond*500) 等价于 time.NewTimer(time.Millisecond*500).C
func main() {
	intChan := make(chan int, 1)
	go func() {
		time.Sleep(time.Second)
		intChan <- 1
	}()

	select {
	case e := <- intChan:
		fmt.Printf("receive: %v\n", e)
	case <-time.NewTimer(time.Millisecond*500).C:
		fmt.Println("timeout")
	//case <-time.After(time.Millisecond*500):
	//	fmt.Println("timeout")
	}
}

// Output:
//timeout

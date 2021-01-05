package main

import (
	"fmt"
)

// close 函数不能传入 nil
// 一个 chan 不能重复 close
// 否则都会引发 runtime panic
func main() {
	dataChan := make(chan int, 5)
	syncChan := make(chan struct{}, 1)
	syncChan1 := make(chan struct{}, 2)
	go func() {
		<- syncChan
		for {
			if elem, ok := <- dataChan; ok {
				fmt.Printf("receive: %d [receiver]\n", elem)
			} else {
				break
			}
		}
		fmt.Println("Done. [receiver]")
		syncChan1 <- struct{}{}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			dataChan <- i
			fmt.Printf("sent: %d [sender]\n", i)
		}
		close(dataChan)
		syncChan <- struct{}{}
		fmt.Println("Done. [sender]")
		syncChan1 <- struct{}{}
	}()

	<- syncChan1
	<- syncChan1
}

// Output:
//sent: 0 [sender]
//sent: 1 [sender]
//sent: 2 [sender]
//sent: 3 [sender]
//sent: 4 [sender]
//Done. [sender]
//receive: 0 [receiver]
//receive: 1 [receiver]
//receive: 2 [receiver]
//receive: 3 [receiver]
//receive: 4 [receiver]
//Done. [receiver]

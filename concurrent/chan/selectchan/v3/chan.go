package main

import (
	"fmt"
)

// select 通常会放在一个单独的 goroutine 中，避免阻塞 main goroutine
func main() {
	intChan := make(chan int, 10)
	for i := 0; i < 10; i++ {
		intChan <- i
	}
	close(intChan)
	syncChan := make(chan struct{}, 1)
	go func() {
	Loop:
		for {
			select {
			case elem, ok := <-intChan:
				if !ok {
					fmt.Println("end")
					break Loop // 如果没有 Loop tag，break 只能结束 select 语句
				}
				fmt.Printf("receive: %v\n", elem)
			}
		}
		syncChan <- struct{}{}
	}()

	<-syncChan
}

// Output:
//receive: 0
//receive: 1
//receive: 2
//receive: 3
//receive: 4
//receive: 5
//receive: 6
//receive: 7
//receive: 8
//receive: 9
//end

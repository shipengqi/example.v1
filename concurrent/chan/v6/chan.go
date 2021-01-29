package main

import (
	"fmt"
	"time"
)

// close 函数不能传入 nil
// 一个 chan 不能重复 close
// 否则都会引发 runtime panic
func main() {
	dataChan := make(chan int)
	go func() {
		dataChan <- 1
		close(dataChan)
		fmt.Println("chan closed")
	}()

	fmt.Println(<- dataChan)
	time.Sleep(time.Second)
	fmt.Println("read from closed chan")
	dataChan <- 1
	fmt.Println("read from closed chan")
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

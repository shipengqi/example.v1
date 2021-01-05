package main

import (
	"fmt"
	"time"
)

// range 可以持续从一个通道接收值，直到 chan 关闭。
// 如果 chan 没有元素，或者 chan 为 nil，则会阻塞 goroutine

var dataChan = make(chan string, 5)

func main() {
	syncChan := make(chan struct{}, 1)
	syncChan1 := make(chan struct{}, 2)
	go receive(dataChan, syncChan, syncChan1)
	go send(dataChan, syncChan, syncChan1)
	<- syncChan1
	<- syncChan1
}

func receive(dataChan <-chan string,
	syncChan <-chan struct{},
	syncChan1 chan<- struct{}) {
	<- syncChan
	//for {
	//	if elem, ok := <- dataChan; ok {
	//		fmt.Printf("receive: %s [receiver]\n", elem)
	//	} else {
	//		break
	//	}
	//  fmt.Println("get nothing from dataChan")
	//}
	// 使用 range 子句，就不需要使用上面死循环的方式去迭代 chan，代码更加简洁
	for elem := range dataChan {
		fmt.Printf("receive: %s [receiver]\n", elem)
	}
	fmt.Println("Done. [receiver]")
	syncChan1 <- struct{}{}
}

func send(dataChan chan<- string,
	syncChan chan<- struct{},
	syncChan1 chan<- struct{}) {
	for _, elem := range []string{"a", "b", "c", "d"} {
		dataChan <- elem
		fmt.Printf("sent: %s [sender]\n", elem)
		if elem == "c" {
			syncChan <- struct{}{}
			fmt.Println("sent a sync signal")
		}
	}
	fmt.Println("wait 2 seconds... [sender]")
	time.Sleep(time.Second * 2)
	close(dataChan)
	syncChan1 <- struct{}{}
}

// Output:
//sent: a [sender]
//sent: b [sender]
//sent: c [sender]
//sent a sync signal
//sent: d [sender]
//wait 2 seconds... [sender]
//receive: a [receiver]
//receive: b [receiver]
//receive: c [receiver]
//receive: d [receiver]
//Done. [receiver]

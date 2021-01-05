package main

import (
	"fmt"
	"time"
)

var dataChan = make(chan string, 5)

// 单向 chan 只是对 chan 的使用方式的约束，如 os/signal.Notify 的函数声明：
// func Notify(c chan<- os.Signal, sig ...os.Signal)
// 第一个参数是发送 chan 类型。但是这并不意味着传入一个只能发送的单向 chan，再调用时，还是要传入一个双向 chan。
// 编译器会负责把双向 chan 换成单向 chan。这就是一种使用约束，在函数内只能对 c chan 进行发送操作。

// 双向 chan 可以转换成单向 chan，但是单向 chan 不能转换为双向 chan
// chan 的传递方向是类型的一部分，不同传递方向 chan 就是不同的类型
func main() {
	syncChan := make(chan struct{}, 1)
	syncChan1 := make(chan struct{}, 2)
	go receive(dataChan, syncChan, syncChan1)
	go send(dataChan, syncChan, syncChan1)
	<- syncChan1
	<- syncChan1
}

// receive 函数只能对 dataChan 和 syncChan 进行接收操作
// 只能对 syncChan1 进行发送操作
func receive(dataChan <-chan string,
	syncChan <-chan struct{},
	syncChan1 chan<- struct{}) {
	<- syncChan
	for {
		if elem, ok := <- dataChan; ok {
			fmt.Printf("receive: %s [receiver]\n", elem)
		} else {
			break
		}
	}
	fmt.Println("Done. [receiver]")
	syncChan1 <- struct{}{}
}

// send 函数只能对 dataChan 和 syncChan 进行发送操作
// 只能对 syncChan1 进行发送操作
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

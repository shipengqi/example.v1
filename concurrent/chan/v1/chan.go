package main

import (
	"fmt"
	"time"
)

var c = make(chan string, 3)

func main() {
	// struct{} 是空结构体类型，在 Go 中，空结构体类型是不占用内存的，
	// 并且所有的该类型的变量都拥有相同的内存地址。
	// 用于传递信号的通道，都应该以此类型为元素类型，除非有更多的信息要传递。
	c1 := make(chan struct{}, 1)
	c2 := make(chan struct{}, 2)
	go func() {
		<- c1 // c1 中没有数据，会阻塞 goroutine，直到 c1 的发送方发送数据
		fmt.Println("received a sync signal and wait a second ... [receiver]")
		time.Sleep(time.Second)
		for {
			if elem, ok := <- c; ok {
				fmt.Println("received: ", elem, " [receiver]")
			} else {
				break
			}
			fmt.Println("get nothing")
		}
		fmt.Println("stopped. [receiver]")
		c2 <- struct{}{}
	}()

	go func() {
		for _, elem := range []string{"a", "b", "c", "d"} {
			c <- elem // 这里发送三个 elem 之后，c 通道就满了，并且 c 的接收方因为 c1 没有数据而阻塞
			fmt.Println("sent:", elem, "[sender]")
			if elem == "c" {
				c1 <- struct{}{} // 发送到 c1，这个 c1 的接收方会恢复执行
				fmt.Println("sent a sync signal. [sender]")
			}
			// c1 的接收方会恢复执行，c 的接收方接收数据，向 c 发送第四个 elem
		}
		fmt.Println("wait 2 seconds... [sender]")
		time.Sleep(time.Second * 2)
		// 等待 2 秒，c 的接收方会收到 4 个数据，然后等待片刻，c 被关闭，继续执行
		close(c)
		c2 <- struct{}{}
	}()

	<- c2
	<- c2
}

// Output:
//sent: a [sender]
//sent: b [sender]
//sent: c [sender]
//sent a sync signal. [sender]
//received a sync signal and wait a second ... [receiver]
//received:  a  [receiver]
//received:  b  [receiver]
//received:  c  [receiver]
//received:  d  [receiver]
//sent: d [sender]
//wait 2 seconds... [sender]
//stopped. [receiver]

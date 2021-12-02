package main

import (
	"fmt"
	"time"
)

// 通过 chan+select 的方式控制多个 goroutine
// 这种方式有很大的局限性
// 如果有很多 goroutine 都需要控制结束怎么办？
// 如果这些 goroutine 又衍生了其他更多的 goroutine 怎么办？
// 这就非常复杂了，即使我们定义很多 chan 也很难解决这个问题，因为 goroutine 的关系链就导致了这种场景非常复杂。
// 下面的示例只是控制 3 个 goroutine，通知 3 个 goroutine 结束，需要分别对 3 个 goroutine 发送信号，也就是要发送三次。
// 或者针对每个 goroutine 初始化一个 chan，分别发送结束信号。
// 这就是 chan+select 的局限性

var stop = make(chan struct{})

func main() {
	go watch("[monitor 1]")
	go watch("[monitor 2]")
	go watch("[monitor 3]")

	time.Sleep(10 * time.Second)
	fmt.Println("stop monitor")
	stop <- struct{}{}
	// stop <- struct{}{}
	// stop <- struct{}{}
	time.Sleep(5 * time.Second)
}

func watch(name string) {
	for {
		select {
		case <-stop:
			fmt.Println(name, "monitor end ...")
			return
		default:
			fmt.Println(name, "goroutine monitoring ...")
			time.Sleep(2 * time.Second)
		}
	}
}

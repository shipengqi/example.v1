package main

import (
	"context"
	"fmt"
	"time"
)

// Context 控制多个 goroutine，就可以解决 chan+select 方式的局限性
// 子 ctx 的 CancelFunc 不能撤销父 ctx
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(ctx)
	go watch(ctx, "[monitor 1]")
	go watch(ctx, "[monitor 2]")
	go watch(ctx, "[monitor 3]")
	go watch(ctx2, "[child monitor 1]")
	time.Sleep(5 * time.Second)
	fmt.Println("cancel child ctx")
	cancel2()
	time.Sleep(5 * time.Second)
	fmt.Println("cancel parent ctx")
	cancel()
	fmt.Println("end")
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "monitor end ...")
			return
		default:
			fmt.Println(name, "goroutine monitoring ...")
			time.Sleep(1 * time.Second)
		}
	}
}

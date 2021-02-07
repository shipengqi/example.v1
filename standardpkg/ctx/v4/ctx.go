package main

import (
	"context"
	"fmt"
	"time"
)

// Context 控制多个 goroutine，就可以解决 chan+select 方式的局限性
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx, "[monitor 1]")
	go watch(ctx, "[monitor 2]")
	go watch(ctx, "[monitor 3]")

	time.Sleep(10 * time.Second)
	fmt.Println("stop monitor")
	cancel()
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
			time.Sleep(2 * time.Second)
		}
	}
}

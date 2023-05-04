package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("monitor end...")
				return
			default:
				fmt.Println("goroutine monitoring ...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("monitor 2 end...")
				return
			default:
				fmt.Println("goroutine 2 monitoring ...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(10 * time.Second)
	fmt.Println("stop monitor")
	cancel()
	time.Sleep(5 * time.Second)

}

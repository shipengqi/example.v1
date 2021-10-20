package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main()  {
	go sigReceiver("receiver1")
	go sigReceiver("receiver2")
	fmt.Println("block ...")
	time.Sleep(time.Second*10)
}


func sigReceiver(flag string) {
	var once sync.Once
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		fmt.Println("Receiving ...")
		select {
		case sig := <-quit:
			once.Do(func() {
				fmt.Printf("%s get a signal %s\n", flag, sig.String())
			})
			fmt.Println("Received ...")
			return
		}
	}
}

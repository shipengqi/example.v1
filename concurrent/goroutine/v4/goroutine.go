package main

import (
	"fmt"
	"time"
)

func main() {
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		go fmt.Printf("v1 Hello, %s!\n", name)
	}
	time.Sleep(time.Second)
	for _, name := range names {
		name := name
		go func() {
			fmt.Printf("v2 Hello, %s!\n", name)
		}()
	}
	time.Sleep(time.Second)
	for _, name := range names {
		go func() {
			fmt.Printf("v3 Hello, %s!\n", name)
		}()
	}
	time.Sleep(time.Second)
	for _, name := range names {
		go func(name string) {
			fmt.Printf("v4 Hello, %s!\n", name)
		}(name)
	}
	time.Sleep(time.Second)
}

// 输出顺序不是一定的
//v1 Hello, Eric!
//v1 Hello, Mark!
//v1 Hello, Robert!
//v1 Hello, Jim!
//v1 Hello, Harry!
//v2 Hello, Mark!
//v2 Hello, Robert!
//v2 Hello, Jim!
//v2 Hello, Eric!
//v2 Hello, Harry!
//v3 Hello, Mark!
//v3 Hello, Mark!
//v3 Hello, Mark!
//v3 Hello, Mark!
//v3 Hello, Mark!
//v4 Hello, Mark!
//v4 Hello, Eric!
//v4 Hello, Harry!
//v4 Hello, Robert!
//v4 Hello, Jim!

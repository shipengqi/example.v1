package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	fmt.Println("A")
	m.Lock()

	go func() {
		time.Sleep(2 *time.Second)
		m.Unlock()
	}()

	// If the lock is already in use, the calling goroutine
	// blocks until the mutex is available.
	m.Lock()
	fmt.Println("B")
}

// Output:
// A
// B

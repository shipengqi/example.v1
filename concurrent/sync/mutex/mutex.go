package main

import (
	"fmt"
	"sync"
	"time"
)

// 注意：对一个未锁定的互斥锁进行解锁操作，会引发 panic，所以最好使用 defer 解锁
// go 1.8 以前是可以恢复运行时 panic 的。但是会导致严重的问题，比如 重复解锁的 goroutine 可能会被永久阻塞
// 所以要尽量避免这种运行时 panic

// 同一个互斥锁的锁定和解锁应该放在同一个层次的代码块中，例如在同一个函数中锁定和解锁，也可以把互斥锁作为一个结构体字段
// 避免在不相关的流程中被误用，导致 panic
// 互斥锁的作用域要尽量小

func main()  {
	var mutex sync.Mutex
	fmt.Println("lock the func main")
	mutex.Lock()
	fmt.Println("locked main")
   for i := 0; i < 3; i ++ {
   	go func(i int) {
   		fmt.Printf("lock the (g%d)\n", i)
   		// 这里的 lock 会阻塞当前的 goroutine，因为使用的是同一个互斥锁 mutex
   		// 一个已经锁定的互斥锁，再其他 goroutine 中重复锁定，就会被阻塞，直到该互斥锁被解锁
   		mutex.Lock()
			fmt.Printf("locked the (g%d)\n", i)
		}(i)
	}
	time.Sleep(time.Second)
   fmt.Println("unlock the func main")
   mutex.Unlock()
   fmt.Println("unlocked main")
   time.Sleep(time.Second)
	repeatUnlock()
}

// 输出：
//lock the func main
//locked main
//lock the (g2)
//lock the (g0)
//lock the (g1)
//unlock the func main
//unlocked main
//locked the (g2)

func repeatUnlock() {
	defer func() {
		fmt.Println("try to recover the panic")
		if p := recover(); p != nil {
			fmt.Printf("recovered panic(%#v)\n", p)
		}
	}()
	var mutex sync.Mutex
	fmt.Println("lock the func repeatedUnlock")
	mutex.Lock()
	fmt.Println("locked repeatedUnlock")
	fmt.Println("unlock the func repeatedUnlock")
	mutex.Unlock()
	fmt.Println("unlocked repeatedUnlock")
	fmt.Println("unlock the func repeatedUnlock again")
	// 无法捕获这个 panic
	mutex.Unlock()
	// panic(errors.New("failed"))
}
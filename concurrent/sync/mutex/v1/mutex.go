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

func main() {
	var mutex sync.Mutex
	fmt.Println("lock the func main")
	mutex.Lock()
	fmt.Println("locked main")
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Printf("try to lock the (g%d)\n", i)
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
// lock the func main
// locked main
// try to lock the (g2)
// try to lock the (g0)
// try to lock the (g1)
// unlock the func main
// unlocked main
// locked the (g2)
// lock the func repeatedUnlock
// locked repeatedUnlock
// unlock the func repeatedUnlock
// unlocked repeatedUnlock
// unlock the func repeatedUnlock again
// fatal error: sync: unlock of unlocked mutex
//
// goroutine 1 [running]:
// runtime.throw({0xa1d335, 0x1})
//	C:/Program Files/Go/src/runtime/panic.go:1198 +0x76 fp=0xc00007de08 sp=0xc00007ddd8 pc=0x9a3116
// sync.throw({0xa1d335, 0x0})
//	C:/Program Files/Go/src/runtime/panic.go:1184 +0x1e fp=0xc00007de28 sp=0xc00007de08 pc=0x9c80fe
// sync.(*Mutex).unlockSlow(0xc00000a0e0, 0xffffffff)
//	C:/Program Files/Go/src/sync/mutex.go:196 +0x3c fp=0xc00007de50 sp=0xc00007de28 pc=0x9da17c
// sync.(*Mutex).Unlock(...)
//	C:/Program Files/Go/src/sync/mutex.go:190
// main.repeatUnlock()
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:65 +0x205 fp=0xc00007def0 sp=0xc00007de50 pc=0x9fe565
// main.main()
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:36 +0x236 fp=0xc00007df80 sp=0xc00007def0 pc=0x9fe1d6
// runtime.main()
//	C:/Program Files/Go/src/runtime/proc.go:255 +0x217 fp=0xc00007dfe0 sp=0xc00007df80 pc=0x9a56d7
// runtime.goexit()
//	C:/Program Files/Go/src/runtime/asm_amd64.s:1581 +0x1 fp=0xc00007dfe8 sp=0xc00007dfe0 pc=0x9ccc81
//
// goroutine 6 [semacquire]:
// sync.runtime_SemacquireMutex(0x16, 0xa8, 0x0)
//	C:/Program Files/Go/src/runtime/sema.go:71 +0x25
// sync.(*Mutex).lockSlow(0xc00000a0b8)
//	C:/Program Files/Go/src/sync/mutex.go:138 +0x165
// sync.(*Mutex).Lock(...)
//	C:/Program Files/Go/src/sync/mutex.go:81
// main.main.func1(0x0)
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:27 +0x8f
// created by main.main
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:23 +0xe5
//
// goroutine 7 [semacquire]:
// sync.runtime_SemacquireMutex(0x16, 0xa8, 0x0)
//	C:/Program Files/Go/src/runtime/sema.go:71 +0x25
// sync.(*Mutex).lockSlow(0xc00000a0b8)
//	C:/Program Files/Go/src/sync/mutex.go:138 +0x165
// sync.(*Mutex).Lock(...)
//	C:/Program Files/Go/src/sync/mutex.go:81
// main.main.func1(0x0)
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:27 +0x8f
// created by main.main
//	C:/Code/example.v1/concurrent/sync/mutex/mutex.go:23 +0xe5

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

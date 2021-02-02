package main

import (
	"fmt"
	"sync"
	"time"
)

// 读写锁就是可以分别对读操作和写操作进行锁定和解锁
// 读-读 不存在互斥
// 读-写 互斥
// 写-写 互斥

func main() {
	var rwm sync.RWMutex
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Printf("try to lock for reading g[%d]\n", i)
			rwm.RLock()
			fmt.Printf("locked for reading g[%d]\n", i)
			time.Sleep(time.Second * 2)
			fmt.Printf("try to unlock for reading g[%d]\n", i)
			rwm.RUnlock()
			fmt.Printf("unlocked for reading g[%d]\n", i)
		}(i)
	}
	time.Sleep(time.Second * 1)
	fmt.Println("try to lock for writing")
	rwm.Lock()
	fmt.Println("locked for writing")
}

// Output:
//try to lock for reading g[2]
//locked for reading g[2]
//try to lock for reading g[1]
//locked for reading g[1]
//try to lock for reading g[0]
//locked for reading g[0]
//try to lock for writing
//try to unlock for reading g[2]
//unlocked for reading g[2]
//try to unlock for reading g[1]
//unlocked for reading g[1]
//try to unlock for reading g[0]
//unlocked for reading g[0]
//locked for writing
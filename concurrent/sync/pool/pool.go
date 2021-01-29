package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// sync.Pool 是一个可以存或取的临时对象集合
// sync.Pool 可以安全被多个线程同时使用，保证线程安全
// sync.Pool 中保存的任何项都可能随时不做通知的释放掉，所以不适合用于像 socket 长连接或数据库连接池。
// sync.Pool 主要用途是增加临时对象的重用率，减少 GC 负担

func main() {
	// 禁用 GC，并在执行结束后恢复
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	var count int32
	pool := sync.Pool{New: func() interface{} {
		return atomic.AddInt32(&count, 1)
	}}

	v1 := pool.Get()
	fmt.Printf("Value 1: %v\n", v1)

	pool.Put(10)
	pool.Put(11)
	pool.Put(12)

	v2 := pool.Get()
	fmt.Printf("Value 2: %v\n", v2)

	debug.SetGCPercent(100)
	runtime.GC() // 主动触发 GC， 会把 pool 中的值，全部回收

	v3 := pool.Get()
	fmt.Printf("Value 3: %v\n", v3)

	pool.New = nil
	v4 := pool.Get()
	fmt.Printf("Value 4: %v\n", v4)

}

// Value 1: 1
// Value 2: 10
// Value 3: 2
// Value 4: <nil>


// sync.Pool 源码中的 init 函数：
// func init() {
//	 runtime_registerPoolCleanup(poolCleanup)
// }
// init 的时候注册了一个 PoolCleanup 函数，负责清除掉 sync.Pool 中的所有的缓存的对象，这个注册函数会在每次 GC 的时候运行
// 所以 sync.Pool 中的值只在两次 GC 中间的时段有效
// 所以上面的运行结果也可能是：

// Value 1: 1
// Value 2: 10
// Value 3: 11
// Value 4: 12

// 所以不要对 pool 中的值有任何假设，它取出的值可能是池中的任何一个值，也可能是 New 函数新生成的
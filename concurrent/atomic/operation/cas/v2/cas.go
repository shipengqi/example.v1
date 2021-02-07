package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/shipengqi/example.v1/utils"
)

type V struct {
	version uint32
	v       int32
}

var (
	src V
	wg  sync.WaitGroup
)

func main() {

	threadNum := 10

	wg.Add(threadNum)
	fmt.Println("====== old value =======")
	fmt.Println("old V.version", src.version)
	fmt.Println("old V.v", src.v)

	for i := 0; i < threadNum; i++ {
		go addValue(1)
	}

	wg.Wait()
	fmt.Println("====== cas value =======")
	fmt.Println("cas V.version", src.version)
	fmt.Println("cas V.v", src.v)
}

// 解决 ABA 问题
func addValue(v int32) {
	defer wg.Done()

	spinNum := 0
	// 在高并发的情况下，单次 CAS 的执行成功率会降低，因此需要配合循环语句 for，形成一个 for+atomic 的类似自旋乐观锁
	for {
		oldV := atomic.LoadInt32(&src.v)
		oldVersion := atomic.LoadUint32(&src.version)
		if atomic.CompareAndSwapInt32(&src.v, oldV, oldV+v) && atomic.CompareAndSwapUint32(&src.version, oldVersion, oldVersion+1) {
			fmt.Println("add value v", src.v)
			fmt.Println("add value version", src.version)
			break
		} else {
			spinNum++
		}
	}
	fmt.Printf("thread: %d,spinnum: %d\n", utils.GoID(), spinNum)
}

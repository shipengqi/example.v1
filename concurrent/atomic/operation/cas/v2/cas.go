package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type V struct {
	version int32
	v       int32
}

var value atomic.Value

func main()  {
	src := V{
		version: 0,
		v:       0,
	}
	value.Store(src)
	fmt.Println("====== old value =======")
	old := value.Load()
	if o, ok := old.(V); ok {
		fmt.Println("old value", o.version, o.v)
	}
	addValue(3)
	fmt.Println("====== cas value =======")
	fmt.Println(value)
}

func addValue(v int32)  {

	// 在高并发的情况下，单次 CAS 的执行成功率会降低，因此需要配合循环语句 for，形成一个 for+atomic 的类似自旋乐观锁
	for {
		old := value.Load()
		if o, ok := old.(V); ok {
			if atomic.CompareAndSwapInt32(&o.version, o.version, int32(time.Now().Nanosecond())) {
				if atomic.CompareAndSwapInt32(&o.v, o.v, o.v + v) {
					break
				}
			}
		}
	}

}
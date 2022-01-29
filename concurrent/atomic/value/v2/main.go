package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// atomic.Value, 此类型相当于一个容器，被用来"原子地"存储（Store）和加载任意类型的值
// 原子性的读取任意结构操作
// func (v *Value) Load() (x interface{}) {
//    // 将*Value指针类型转换为*ifaceWords指针类型
//	vp := (*ifaceWords)(unsafe.Pointer(v))
//	// 原子性的获取到v的类型typ的指针
//	typ := LoadPointer(&vp.typ)
//	// 如果没有写入或者正在写入，先返回，^uintptr(0)代表过渡状态，见下文
//	if typ == nil || uintptr(typ) == ^uintptr(0) {
//		return nil
//	}
//	// 原子性的获取到v的真正的值data的指针，然后返回
//	data := LoadPointer(&vp.data)
//	xp := (*ifaceWords)(unsafe.Pointer(&x))
//	xp.typ = typ
//	xp.data = data
//	return
// }
//
// 原子性的存储任意结构操作
// func (v *Value) Store(x interface{}) {
//	if x == nil {
//		panic("sync/atomic: store of nil value into Value")
//	}
//	// 将现有的值和要写入的值转换为ifaceWords类型，这样下一步就能获取到它们的原始类型和真正的值
//	vp := (*ifaceWords)(unsafe.Pointer(v))
//	xp := (*ifaceWords)(unsafe.Pointer(&x))
//	for {
//		// 获取现有的值的type
//		typ := LoadPointer(&vp.typ)
//		// 如果typ为nil说明这是第一次Store
//		if typ == nil {
//			// 如果你是第一次，就死死占住当前的processor，不允许其他goroutine再抢
//			runtime_procPin()
//			// 使用 CAS 操作，先尝试将 typ 设置为 ^uintptr(0) 这个中间状态
//			// 如果失败，则证明已经有别的线程抢先完成了赋值操作
//			// 那它就解除抢占锁，然后重新回到 for 循环第一步
//			if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
//				runtime_procUnpin()
//				continue
//			}
//			// 如果设置成功，说明当前goroutine中了jackpot
//			// 那么就原子性的更新对应的指针，最后解除抢占锁
//			StorePointer(&vp.data, xp.data)
//			StorePointer(&vp.typ, xp.typ)
//			runtime_procUnpin()
//			return
//		}
//		// 如果typ为^uintptr(0)说明第一次写入还没有完成，继续循环等待
//		if uintptr(typ) == ^uintptr(0) {
//			continue
//		}
//		// 如果要写入的类型和现有的类型不一致，则panic
//		if typ != xp.typ {
//			panic("sync/atomic: store of inconsistently typed value into Value")
//		}
//		// 更新data
//		StorePointer(&vp.data, xp.data)
//		return
//	}
// }

func main() {
	config := atomic.Value{}
	config.Store(22)

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			// 在某一个goroutine中修改配置
			if i == 0 {
				config.Store(23)
			}
			// 输出中夹杂 22，23
			fmt.Println(i, config.Load())
		}(i)
	}
	wg.Wait()
}

// 9 22
// 4 22
// 2 23
// 0 23
// 1 23
// 3 22
// 5 23
// 6 23
// 7 23
// 8 23

package main

import (
	"fmt"
	"sync/atomic"
)
// 原子操作中最经典的 CAS(compare-and-swap) 函数。
//
// - func CompareAndSwapInt32(addr \*int32, old, new int32) (swapped bool)
// - func CompareAndSwapInt64(addr \*int64, old, new int64) (swapped bool)
// - func CompareAndSwapPointer(addr \*unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
// - func CompareAndSwapUint32(addr \*uint32, old, new uint32) (swapped bool)
// - func CompareAndSwapUint64(addr \*uint64, old, new uint64) (swapped bool)
// - func CompareAndSwapUintptr(addr \*uintptr, old, new uintptr) (swapped bool)
//
// CAS 的意思是判断内存中的某个值是否等于 old 值，如果是的话，则赋 new 值给这块内存。CAS 是一个方法，并不局限在 CPU 原子操作中。
// CAS 是乐观锁，但是也就代表 CAS 是有赋值不成功的时候，调用 CAS 的那一方就需要处理赋值不成功的后续行为了。
//
// 这一系列的函数需要比较后再进行交换，也有不需要进行比较就进行交换的原子操作。
//
// - func SwapInt32(addr \*int32, new int32) (old int32)
// - func SwapInt64(addr \*int64, new int64) (old int64)
// - func SwapPointer(addr \*unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
// - func SwapUint32(addr \*uint32, new uint32) (old uint32)
// - func SwapUint64(addr \*uint64, new uint64) (old uint64)
// - func SwapUintptr(addr \*uintptr, new uintptr) (old uintptr)

var value int32

func main()  {
	fmt.Println("====== old value =======")
	fmt.Println(value)
	addValue(3)
	fmt.Println("====== cas value =======")
	fmt.Println(value)
}

func addValue(v int32)  {

	// 在高并发的情况下，单次 CAS 的执行成功率会降低，因此需要配合循环语句 for，形成一个 for+atomic 的类似自旋乐观锁
	for {
		// 在进行读取 value 的操作的过程中,其他对此值的读写操作是可以被同时进行的,那么这个读操作很可能会读取到一个只被修改了一半的数据.
		// 因此要使用原子读取
		old := atomic.LoadInt32(&value)
		if atomic.CompareAndSwapInt32(&value, old, old + v) {
			break
		}
	}

}
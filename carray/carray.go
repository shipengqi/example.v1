package carray

import (
	"errors"
	"fmt"
	"sync/atomic"
)

type ConcurrentArray interface {
	Set(index uint32, elem int) error
	Get(index uint32) (int, error)
	Len() uint32
}

type intArray struct {
	length uint32
	val atomic.Value
}

func NewConcurrentArray(len uint32) ConcurrentArray {
	nativeArr := make([]int, len)
	arr := &intArray{length: len}
	arr.val.Store(nativeArr)
	return arr
}

func (i *intArray) Len() uint32 {
	return i.length
}

func (i *intArray) Set(index uint32, elem int) (err error) {
	if err = i.checkIndex(index); err != nil {
		return
	}
	if err = i.checkValue(); err != nil {
		return
	}

	// 不要这样做！否则会形成竞态条件！
	// 因为切片是引用类型，所以 oldArray[index] = elem ，相当于直接操作原切片
	// 并发调用时，会产生竞态条件
	// 如果这里是数组就不会有问题
	// oldArray := array.val.Load().([]int)
	// oldArray[index] = elem
	// array.val.Store(oldArray) // 这行代码没有任何作用


	newArray := make([]int, i.length)
	copy(newArray, i.val.Load().([]int))
	newArray[index] = elem
	i.val.Store(newArray)
	return
}

func (i *intArray) Get(index uint32) (elem int, err error) {
	if err = i.checkIndex(index); err != nil {
		return
	}
	if err = i.checkValue(); err != nil {
		return
	}
	elem = i.val.Load().([]int)[index]
	return
}

func (i *intArray) checkIndex(index uint32) error {
	if index >= i.length {
		return fmt.Errorf("Index out of range [0, %d) ", i.length)
	}
	return nil
}

func (i *intArray) checkValue() error {
	v := i.val.Load()
	if v == nil {
		return errors.New("Invalid int array ")
	}
	return nil
}
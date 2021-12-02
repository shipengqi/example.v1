package cmap

import (
	"bytes"
	"sync"
	"sync/atomic"
)

type Bucket interface {
	// Put 会放入一个键-元素对
	// 第一个返回值表示是否新增了键-元素对
	// 若在调用此方法前已经锁定 lock，则不要把 lock 传入！否则必须传入对应的 lock
	Put(p Pair, lock sync.Locker) (bool, error)
	// Get 会获取指定键的键-元素对。
	Get(key string) Pair
	// GetFirstPair 会返回第一个键-元素对
	GetFirstPair() Pair
	// Delete 会删除指定的键-元素对
	// 若在调用此方法前已经锁定 lock，则不要把 lock 传入！否则必须传入对应的 lock
	Delete(key string, lock sync.Locker) bool
	// Clear 会清空当前散列桶。
	// 若在调用此方法前已经锁定 lock，则不要把 lock 传入！否则必须传入对应的 lock
	Clear(lock sync.Locker)
	// Size 会返回当前散列桶的尺寸
	Size() uint64
	// String 会返回当前散列桶的字符串表示形式
	String() string
}

// 使用单链表和原子操作消除了 读-读 读-写 操作之间的竞态条件
// 这种无锁化实现，大大提高了读操作的性能。但是写操作还是需要使用互斥锁来消除
type bucket struct {
	// firstValue 存储的是键-元素对列表的表头。
	firstValue atomic.Value
	size       uint64
}

// 占位符。
// 由于原子值不能存储 nil，所以当散列桶空时用此符占位
var placeholder Pair = &pair{}

func newBucket() Bucket {
	b := &bucket{}
	b.firstValue.Store(placeholder)
	return b
}

func (b *bucket) Put(p Pair, lock sync.Locker) (bool, error) {
	if p == nil {
		return false, newParamError("pair is nil")
	}
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		b.firstValue.Store(p)
		atomic.AddUint64(&b.size, 1)
		return true, nil
	}
	var target Pair
	key := p.Key()
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			break
		}
	}
	if target != nil {
		_ = target.SetElement(p.Element())
		return false, nil
	}
	// 无锁化实现，firstValue 的操作总是原子的。key value 对的添加，只会把 p 变为新的表头，并把 next 指向原来的表头。
	// 新表头后面的链表每个节点都没有变化，这样就可以保证 key value 对的获取操作，任何时候都可以原子的获取一个表头，并可以并发安全的遍历链表
	_ = p.SetNext(firstPair)
	b.firstValue.Store(p)
	atomic.AddUint64(&b.size, 1)
	return true, nil
}

func (b *bucket) Get(key string) Pair {
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return nil
	}
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			return v
		}
	}
	return nil
}

func (b *bucket) GetFirstPair() Pair {
	if v := b.firstValue.Load(); v == nil {
		return nil
	} else if p, ok := v.(Pair); !ok || p == placeholder {
		return nil
	} else {
		return p
	}
}

func (b *bucket) Delete(key string, lock sync.Locker) bool {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	firstPair := b.GetFirstPair()
	if firstPair == nil {
		return false
	}
	var prevPairs []Pair // 副本
	var target Pair
	var breakpoint Pair
	for v := firstPair; v != nil; v = v.Next() {
		if v.Key() == key {
			target = v
			breakpoint = v.Next()
			break
		}
		prevPairs = append(prevPairs, v)
	}
	if target == nil {
		return false
	}
	newFirstPair := breakpoint
	for i := len(prevPairs) - 1; i >= 0; i-- {
		pairCopy := prevPairs[i].Copy()
		pairCopy.SetNext(newFirstPair)
		newFirstPair = pairCopy
	}
	if newFirstPair != nil {
		// 无锁化实现，原子替换表头前，任何读取操作都会访问旧的单链表，表头替换完成，所有读取操作就会访问新的链表
		// 即使有先开始遍历的读取操作，也不会受到影响。所以是并发安全的
		b.firstValue.Store(newFirstPair)
	} else {
		b.firstValue.Store(placeholder)
	}
	atomic.AddUint64(&b.size, ^uint64(0))
	return true
}

func (b *bucket) Clear(lock sync.Locker) {
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	}
	atomic.StoreUint64(&b.size, 0)
	b.firstValue.Store(placeholder)
}

func (b *bucket) Size() uint64 {
	return atomic.LoadUint64(&b.size)
}

func (b *bucket) String() string {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	for v := b.GetFirstPair(); v != nil; v = v.Next() {
		buf.WriteString(v.String())
		buf.WriteString(" ")
	}
	buf.WriteString("]")
	return buf.String()
}

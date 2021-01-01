package cmap

import (
	"math"
	"sync/atomic"
)

type ConcurrentMap interface {
	// 返回并发量
	Concurrency() int
	Put(key string, element interface{}) (bool, error)
	Get(key string) interface{}
	Delete(key string) bool
	Len() uint64
}

// concurrency 并发量，也代表了 segments 的长度
// 散列段，每个散列段要提供对其包含的键值对的读写操作，通过锁来保证并发安全性。
// 多个散列段就由多个互斥锁保护。这样的加锁方式叫做 分段锁。分段锁可以降低互斥锁带来的开销，同时保护资源。
// 因为同一时刻同一个散列段中的键值对只能被一个 goroutine 进行读写。但是不同的散列段可以并发访问。
// 如果 concurrency 为 16，那么就可以有 16 个 goroutine 并发访问 ConcurrentMap。
type ConcurrentMapImpl struct {
	concurrency int
	segments    []Segment
	total       uint64
}

func NewConcurrentMap(concurrency int, pairRedistributor PairRedistributor) (ConcurrentMap, error) {
	if concurrency <= 0 {
		return nil, newParamError("concurrency is too small")
	}
	if concurrency > MAX_CONCURRENCY {
		return nil, newParamError("concurrency is too large")
	}
	segments := make([]Segment, concurrency)
	for i := 0; i < concurrency; i++ {
		segments[i] = newSegment(DEFAULT_BUCKET_NUMBER, pairRedistributor)
	}
	return &ConcurrentMapImpl{
		concurrency: concurrency,
		segments:    segments,
		total:       0,
	}, nil
}


func (c *ConcurrentMapImpl) Concurrency() int {
	return c.concurrency
}

func (c *ConcurrentMapImpl) Put(key string, element interface{}) (bool, error) {
	p, err := newPair(key, element)
	if err != nil {
		return false, err
	}
	s := c.findSegment(p.Hash())
	ok, err := s.Put(p)
	if ok {
		atomic.AddUint64(&c.total, 1)
	}
	return ok, err
}

func (c *ConcurrentMapImpl) Get(key string) interface{} {
	keyHash := hash(key)
	s := c.findSegment(keyHash)
	pair := s.GetWithHash(key, keyHash)
	if pair == nil {
		return nil
	}
	return pair.Element()
}

func (c *ConcurrentMapImpl) Delete(key string) bool {
	s := c.findSegment(hash(key))
	if s.Delete(key) {
		atomic.AddUint64(&c.total, ^uint64(0))
		return true
	}
	return false
}

func (c *ConcurrentMapImpl) Len() uint64 {
	return atomic.LoadUint64(&c.total)
}

func (c *ConcurrentMapImpl) findSegment(keyHash uint64) Segment {
	if c.concurrency == 1 {
		return c.segments[0]
	}
	var keyHashHigh int
	if keyHash > math.MaxUint32 {
		keyHashHigh = int(keyHash >> 48)
	} else {
		keyHashHigh = int(keyHash >> 16)
	}
	return c.segments[keyHashHigh%c.concurrency]
}

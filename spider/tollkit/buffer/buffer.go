package buffer

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Buffer 代表 FIFO 的缓冲器的接口类型
type Buffer interface {
	// Cap 用于获取本缓冲器的容量
	Cap() uint32
	// Len 用于获取本缓冲器中的数据数量
	Len() uint32
	// Put 用于向缓冲器放入数据
	// 注意！本方法应该是非阻塞的
	// 若缓冲器已关闭则会直接返回非 nil 的错误值
	Put(datum interface{}) (bool, error)
	// Get 用于从缓冲器获取器
	// 注意！本方法应该是非阻塞的
	// 若缓冲器已关闭则会直接返回非 nil 的错误值
	Get() (interface{}, error)
	// Close 用于关闭缓冲器
	// 若缓冲器之前已关闭则返回 false，否则返回 true
	Close() bool
	// Closed 用于判断缓冲器是否已关闭
	Closed() bool
}

type BufferImpl struct {
	// ch 代表存放数据的通道
	ch chan interface{}
	// closed 代表缓冲器的关闭状态：0 - 未关闭；1 - 已关闭
	closed uint32
	// closingLock 代表为了消除因关闭缓冲器而产生的竞态条件的读写锁
	closingLock sync.RWMutex
}

func NewBufferImpl(size uint32) (Buffer, error) {
	if size == 0 {
		errMsg := fmt.Sprintf("illegal size for buffer: %d", size)
		return nil, errors.New(errMsg)
	}
	return &BufferImpl{
		ch: make(chan interface{}, size),
	}, nil
}

func (b *BufferImpl) Put(datum interface{}) (ok bool, err error) {
	b.closingLock.RLock()
	defer b.closingLock.RUnlock()
	if b.Closed() {
		return false, ErrClosedBuffer
	}
	select {
	case b.ch <- datum:
		ok = true
	default:
		ok = false
	}
	return
}

func (b *BufferImpl) Get() (interface{}, error) {
	select {
	case datum, ok := <-b.ch:
		if !ok {
			return nil, ErrClosedBuffer
		}
		return datum, nil
	default:
		return nil, nil
	}
}

func (b *BufferImpl) Close() bool {
	if atomic.CompareAndSwapUint32(&b.closed, 0, 1) {
		b.closingLock.Lock()
		close(b.ch)
		b.closingLock.Unlock()
		return true
	}
	return false
}

func (b *BufferImpl) Closed() bool {
	if atomic.LoadUint32(&b.closed) == 0 {
		return false
	}
	return true
}
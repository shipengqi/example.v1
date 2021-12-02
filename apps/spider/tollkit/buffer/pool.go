package buffer

import (
	"sync"
	"sync/atomic"
)

// Pool 代表数据缓冲池的接口类型。
type Pool interface {
	// BufferCap 用于获取池中缓冲器的统一容量
	BufferCap() uint32
	// MaxBufferNumber 用于获取池中缓冲器的最大数量
	MaxBufferNumber() uint32
	// BufferNumber 用于获取池中缓冲器的数量
	BufferNumber() uint32
	// Total 用于获取缓冲池中数据的总数
	Total() uint64
	// Put 用于向缓冲池放入数据
	// 注意！本方法应该是阻塞的
	// 若缓冲池已关闭则会直接返回非 nil 的错误值
	Put(datum interface{}) error
	// Get 用于从缓冲池获取数据
	// 注意！本方法应该是阻塞的
	// 若缓冲池已关闭则会直接返回非 nil 的错误值
	Get() (datum interface{}, err error)
	// Close 用于关闭缓冲池
	// 若缓冲池之前已关闭则返回 false，否则返回 true
	Close() bool
	// Closed 用于判断缓冲池是否已关闭
	Closed() bool
}

type PoolImpl struct {
	// bufferCap 代表缓冲器的统一容量
	bufferCap uint32
	// maxBufferNumber 代表缓冲器的最大数量
	maxBufferNumber uint32
	// bufferNumber 代表缓冲器的实际数量
	bufferNumber uint32
	// total 代表池中数据的总数
	total uint64
	// bufCh 代表存放缓冲器的通道
	bufCh chan Buffer
	// closed 代表缓冲池的关闭状态：0-未关闭；1-已关闭
	closed uint32
	// lock 代表保护内部共享资源的读写锁
	rwlock sync.RWMutex
}

func (p *PoolImpl) Put(datum interface{}) (err error) {
	if p.Closed() {
		return ErrClosedBufferPool
	}
	var count uint32
	maxCount := p.BufferNumber() * 5
	var ok bool
	for buf := range p.bufCh {
		ok, err = p.putData(buf, datum, &count, maxCount)
		if ok || err != nil {
			break
		}
	}
	return
}

// putData 用于向给定的缓冲器放入数据，并在必要时把缓冲器归还给池。
func (p *PoolImpl) putData(
	buf Buffer, datum interface{}, count *uint32, maxCount uint32) (ok bool, err error) {
	if p.Closed() {
		return false, ErrClosedBufferPool
	}
	defer func() {
		p.rwlock.RLock()
		if p.Closed() {
			atomic.AddUint32(&p.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			p.bufCh <- buf
		}
		p.rwlock.RUnlock()
	}()
	ok, err = buf.Put(datum)
	if ok {
		atomic.AddUint64(&p.total, 1)
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已满而未放入数据就递增计数。
	(*count)++
	// 如果尝试向缓冲器放入数据的失败次数达到阈值，
	// 并且池中缓冲器的数量未达到最大值，
	// 那么就尝试创建一个新的缓冲器，先放入数据再把它放入池。
	if *count >= maxCount &&
		p.BufferNumber() < p.MaxBufferNumber() {
		p.rwlock.Lock()
		if p.BufferNumber() < p.MaxBufferNumber() {
			if p.Closed() {
				p.rwlock.Unlock()
				return
			}
			newBuf, _ := NewBuffer(p.bufferCap)
			newBuf.Put(datum)
			p.bufCh <- newBuf
			atomic.AddUint32(&p.bufferNumber, 1)
			atomic.AddUint64(&p.total, 1)
			ok = true
		}
		p.rwlock.Unlock()
		*count = 0
	}
	return
}

func (p *PoolImpl) Get() (datum interface{}, err error) {
	if p.Closed() {
		return nil, ErrClosedBufferPool
	}
	var count uint32
	maxCount := p.BufferNumber() * 10
	for buf := range p.bufCh {
		datum, err = p.getData(buf, &count, maxCount)
		if datum != nil || err != nil {
			break
		}
	}
	return
}

// getData 用于从给定的缓冲器获取数据，并在必要时把缓冲器归还给池。
func (p *PoolImpl) getData(
	buf Buffer, count *uint32, maxCount uint32) (datum interface{}, err error) {
	if p.Closed() {
		return nil, ErrClosedBufferPool
	}
	defer func() {
		// 如果尝试从缓冲器获取数据的失败次数达到阈值，
		// 同时当前缓冲器已空且池中缓冲器的数量大于1，
		// 那么就直接关掉当前缓冲器，并不归还给池。
		if *count >= maxCount &&
			buf.Len() == 0 &&
			p.BufferNumber() > 1 {
			buf.Close()
			atomic.AddUint32(&p.bufferNumber, ^uint32(0))
			*count = 0
			return
		}
		p.rwlock.RLock()
		if p.Closed() {
			atomic.AddUint32(&p.bufferNumber, ^uint32(0))
			err = ErrClosedBufferPool
		} else {
			p.bufCh <- buf
		}
		p.rwlock.RUnlock()
	}()
	datum, err = buf.Get()
	if datum != nil {
		atomic.AddUint64(&p.total, ^uint64(0))
		return
	}
	if err != nil {
		return
	}
	// 若因缓冲器已空未取出数据就递增计数。
	(*count)++
	return
}

func (p *PoolImpl) Close() bool {
	if !atomic.CompareAndSwapUint32(&p.closed, 0, 1) {
		return false
	}
	p.rwlock.Lock()
	defer p.rwlock.Unlock()
	close(p.bufCh)
	for buf := range p.bufCh {
		buf.Close()
	}
	return true
}

func (p *PoolImpl) Closed() bool {
	if atomic.LoadUint32(&p.closed) == 1 {
		return true
	}
	return false
}

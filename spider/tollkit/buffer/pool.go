package buffer

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

}
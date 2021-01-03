package v1

import (
	"errors"
	"io"
	"os"
	"sync"
)

// Data 代表数据的类型。
type Data []byte

// DataFile 代表数据文件的接口类型。
// RSN 和 WSN 相当于一个计数值，从 1 开始。可以通过这两个方法得到当前已经读取和写入的数据块的数量
type DataFile interface {
	// Read 会读取一个数据块。
	Read() (rsn int64, d Data, err error)
	// Write 会写入一个数据块。
	Write(d Data) (wsn int64, err error)
	// RSN 会获取最后读取的数据块的序列号。
	RSN() int64
	// WSN 会获取最后写入的数据块的序列号。
	WSN() int64
	// DataLen 会获取数据块的长度。
	DataLen() uint32
	// Close 会关闭数据文件。
	Close() error
}

// dataFile 代表数据文件的实现类型。
type dataFile struct {
	f       *os.File     // 文件
	fmutex  sync.RWMutex // 被用于文件的读写锁
	rcond   *sync.Cond   //读操作需要用到的条件变量
	woffset int64        // 写操作需要用到的偏移量
	roffset int64        // 读操作需要用到的偏移量
	wmutex  sync.Mutex   // 写操作需要用到的互斥锁
	rmutex  sync.Mutex   // 读操作需要用到的互斥锁
	dataLen uint32       // 数据块长度
}

// sync.Cond 通过 NewCond 函数创建，参数必须是 sync.Locker 类型，也就是说可以是一个互斥锁，也可以是一个读写锁
// Cond 类型的是三个方法：
// Wait 等待通知，对 Cond 相关联的 Locker 锁自动解锁，并使它所在的 goroutine 阻塞，等待通知，收到通知后，
// 立即唤醒当前 goroutine。
// Signal 单发通知
// Broadcast 广播通知

func NewDataFile(path string, dataLen uint32) (DataFile, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if dataLen == 0 {
		return nil, errors.New("Invalid data length: ")
	}
	df := &dataFile{f: f, dataLen: dataLen}
	df.rcond = sync.NewCond(df.fmutex.RLocker())
	return df, nil
}

func (df *dataFile) Read() (rsn int64, d Data, err error) {
	// 读取并更新读偏移量。
	var offset int64
	df.rmutex.Lock()
	offset = df.roffset
	df.roffset += int64(df.dataLen)
	df.rmutex.Unlock()

	// 读取一个数据块。
	rsn = offset / int64(df.dataLen)
	bytes := make([]byte, df.dataLen)
	df.fmutex.RLock()
	defer df.fmutex.RUnlock()
	for {
		_, err = df.f.ReadAt(bytes, offset)
		if err != nil {
			if err == io.EOF {
				df.rcond.Wait()
				continue
			}
			return
		}
		d = bytes
		return
	}
}

func (df *dataFile) Write(d Data) (wsn int64, err error) {
	// 读取并更新写偏移量。
	var offset int64
	df.wmutex.Lock()
	offset = df.woffset
	df.woffset += int64(df.dataLen)
	df.wmutex.Unlock()

	//写入一个数据块。
	wsn = offset / int64(df.dataLen)
	var bytes []byte
	if len(d) > int(df.dataLen) {
		bytes = d[0:df.dataLen]
	} else {
		bytes = d
	}
	df.fmutex.Lock()
	defer df.fmutex.Unlock()
	_, err = df.f.Write(bytes)
	df.rcond.Signal()
	return
}

func (df *dataFile) RSN() int64 {
	df.rmutex.Lock()
	defer df.rmutex.Unlock()
	return df.roffset / int64(df.dataLen)
}

func (df *dataFile) WSN() int64 {
	df.wmutex.Lock()
	defer df.wmutex.Unlock()
	return df.woffset / int64(df.dataLen)
}

func (df *dataFile) DataLen() uint32 {
	return df.dataLen
}

func (df *dataFile) Close() error {
	if df.f == nil {
		return nil
	}
	return df.f.Close()
}
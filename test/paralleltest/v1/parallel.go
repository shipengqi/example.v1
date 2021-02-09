package v1

import "sync"

var (
	data   = make(map[string]string)
	locker sync.RWMutex
)

// 如果注释掉 WriteToMap 和 ReadFromMap 中 locker 保护的代码，执行测试，加上 `-race` 一定会失败：`go test -v -race`
// ==================
//    TestReadFromMap: testing.go:906: race detected during execution of test
//    TestWriteToMap: testing.go:906: race detected during execution of test
func WriteToMap(k, v string) {
	locker.Lock()
	defer locker.Unlock()
	data[k] = v
}

func ReadFromMap(k string) string {
	locker.RLock()
	defer locker.RUnlock()
	return data[k]
}
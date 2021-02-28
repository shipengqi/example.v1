package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// 利用 runtime.Stack 的堆栈信息。`runtime.Stack(buf []byte, all bool) int` 将当前的堆栈信息写入到一个 slice 中，
// 堆栈的第一行为 `goroutine #### […` (如 goroutine 51 [running]:), 其中 `####` 就是当前的 goroutine ID。
// 获取堆栈信息会影响性能，所以建议在 debug 的时候才用
// 获取 goroutine ID
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

// 第二种获取 goroutine ID 的方式，通过汇编获取 runtime·getg 方法的调用结果，这种 hack 的方式(Go team并不想暴露go id的信息),
// 针对不同的 Go 版本中需要特殊的 hack 手段。如 https://github.com/petermattis/goid
// go runtime 中实现了一个 getg 方法，可以获取当前的 goroutine：
//
// getg() alone returns the current g
//
// 但是内部方法。可以尝试直接修改源代码 src/runtime/runtime2.go 中添加 Goid 函数，将 goid 暴露给应用层：
//
// func Goid() int64 {
//    _g_ := getg()
//    return _g_.goid
// }
//
// 这个方式能解决法 1 的性能问题，但是会导致你的程序只能在修改了源代码的机器上才能编译，没有移植性，并且每次 go 版本升级以后，都需要重新修改源代码，
// 维护成本较高。
//
// petermattis/goid 模仿 runtime.getg 暴露出一个 getg 的方法，这种方式和前面的方式的性能差距非常大，一千倍左右
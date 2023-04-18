package main

import (
	"github.com/sourcegraph/conc"
)

/* conc.WaitGroup 对标准库的 wg.Add 与 wg.Done 函数进行了封装
type WaitGroup struct {
	wg sync.WaitGroup
	pc PanicCatcher
}

// Go spawns a new goroutine in the WaitGroup
func (h *WaitGroup) Go(f func()) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		h.pc.Try(f)
	}()
}

h.pc.Try(f) 中封装了 recover 可以捕获 panic，并输出堆栈信息
func (p *PanicCatcher) Try(f func()) {
	defer p.tryRecover()
	f()
}

*/

func main() {
	var wg conc.WaitGroup
	wg.Go(doSomethingThatMightPanic)

	// 捕获 panic
	// wg.WaitAndRecover()

	// Wait 会重新抛出 panic
	// panics with a nice stacktrace
	wg.Wait()

}

func doSomethingThatMightPanic() {
	panic("test panic")
}

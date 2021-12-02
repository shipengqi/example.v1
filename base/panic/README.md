# panic

panic 只会触发当前 Goroutine 的 defer；
recover 只有在 defer 中调用才会生效；
panic 允许在 defer 中嵌套多次调用；

数据结构 [runtime._panic](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/runtime2.go#L942)

```go
type _panic struct {
	argp      unsafe.Pointer
	arg       interface{}
	link      *_panic
	recovered bool
	aborted   bool
	pc        uintptr
	sp        unsafe.Pointer
	goexit    bool
}
```

argp 是指向 defer 调用时参数的指针；
arg 是调用 panic 时传入的参数；
link 指向了更早调用的 runtime._panic 结构；
recovered 表示当前 runtime._panic 是否被 recover 恢复；
aborted 表示当前的 panic 是否被强行终止；

panic 函数可以被连续多次调用，它们之间通过 link 可以组成链表。

编译器会将关键字 panic 转换成 runtime.gopanic，该函数的执行过程包含以下几个步骤：

创建新的 runtime._panic 并添加到所在 Goroutine 的 _panic 链表的最前面；
在循环中不断从当前 Goroutine 的 _defer 中链表获取 runtime._defer 并调用 runtime.reflectcall 运行延迟调用函数；
调用 runtime.fatalpanic 中止整个程序；

编译器会将关键字 recover 转换成 runtime.gorecover：


编译器会负责做转换关键字的工作；
将 panic 和 recover 分别转换成 runtime.gopanic 和 runtime.gorecover；
将 defer 转换成 runtime.deferproc 函数；
在调用 defer 的函数末尾调用 runtime.deferreturn 函数；
在运行过程中遇到 runtime.gopanic 方法时，会从 Goroutine 的链表依次取出 runtime._defer 结构体并执行；
如果调用延迟执行函数时遇到了 runtime.gorecover 就会将 _panic.recovered 标记成 true 并返回 panic 的参数；
在这次调用结束之后，runtime.gopanic 会从 runtime._defer 结构体中取出程序计数器 pc 和栈指针 sp 并调用 runtime.recovery 函数进行恢复程序；
runtime.recovery 会根据传入的 pc 和 sp 跳转回 runtime.deferproc；
编译器自动生成的代码会发现 runtime.deferproc 的返回值不为 0，这时会跳回 runtime.deferreturn 并恢复到正常的执行流程；
如果没有遇到 runtime.gorecover 就会依次遍历所有的 runtime._defer，并在最后调用 runtime.fatalpanic 中止程序、打印 panic 的参数并返回错误码 2；
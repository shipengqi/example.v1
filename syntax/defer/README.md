# defer

[runtime._defer](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/runtime2.go#L907) 结构
体是延迟调用链表上的一个元素，所有的结构体都会通过 link 字段串联成链表。

```go
type _defer struct {
	siz       int32
	started   bool
	openDefer bool
	sp        uintptr
	pc        uintptr
	fn        *funcval
	_panic    *_panic
	link      *_defer
}
```

siz 是参数和结果的内存大小；
sp 和 pc 分别代表栈指针和调用方的程序计数器；
fn 是 defer 关键字中传入的函数；
_panic 是触发延迟调用的结构体，可能为空；
openDefer 表示当前 defer 是否经过开放编码的优化；


堆分配、栈分配和开放编码是处理 defer 关键字的三种方法，早期的 Go 语言会在堆上分配 runtime._defer 结构体，不过该实现的性能较差，Go 语言
在 1.13 中引入栈上分配的结构体，减少了 30% 的额外开销1，并在 1.14 中引入了基于开放编码的 defer

```go
func (s *state) stmt(n *Node) {
	...
	switch n.Op {
	case ODEFER:
		if s.hasOpenDefers {
			s.openDeferRecord(n.Left) // 开放编码
		} else {
			d := callDefer // 堆分配
			if n.Esc == EscNever {
				d = callDeferStack // 栈分配
			}
			s.callResult(n.Left, d)
		}
	}
}
```

Go 语言的编译器不仅将 defer 转换成了 runtime.deferproc，还在所有调用 defer 的函数结尾插入了 runtime.deferreturn。上述两个运行时函数是 defer 关键字运行时机制的入口，它们分别承担了不同的工作：

runtime.deferproc 负责创建新的延迟调用；
runtime.deferreturn 负责在函数调用结束时执行所有的延迟调用；

runtime.deferproc 中 runtime.newdefer 的作用是想尽办法获得 runtime._defer 结构体，这里包含三种路径：

从调度器的延迟调用缓存池 sched.deferpool 中取出结构体并将该结构体追加到当前 Goroutine 的缓存池中；
从 Goroutine 的延迟调用缓存池 pp.deferpool 中取出结构体；
通过 runtime.mallocgc 在堆上创建一个新的结构体；

只要获取到 runtime._defer 结构体，它都会被追加到所在 Goroutine _defer 链表的最前面。

defer 关键字的插入顺序是从后向前的，而 defer 关键字执行是从前向后的，这也是为什么后调用的 defer 会优先执行。

runtime.deferreturn 会从 Goroutine 的 _defer 链表中取出最前面的 runtime._defer 并调用 runtime.jmpdefer 传入需要执行的函数和参数：

runtime.jmpdefer 是一个用汇编语言实现的运行时函数，它的主要工作是跳转到 defer 所在的代码段并在执行结束之后跳转回 runtime.deferreturn。
# select

`select` 在 Go 语言的源代码中不存在对应的结构体，但是使用 `runtime.scase` 结构体表示 `select` 控制结构中的 `case`：
```go
type scase struct {
    c    *hchan         // chan
    elem unsafe.Pointer // data element
}
```

因为非默认的 case 中都与 Channel 的发送和接收有关，所以 `runtime.scase` 结构体中也包含一个 `runtime.hchan` 类型的字段存储 case 中使用
的 Channel。

`select` 语句在编译期间会被转换成 `OSELECT` 节点。每个 `OSELECT` 节点都会持有一组 `OCASE` 节点，如果 `OCASE` 的执行条件是空，那就意味着这是一
个 `default` 节点。

编译器在中间代码生成期间会根据 `select` 中 case 的不同对控制语句进行优化，[cmd/compile/internal/gc.walkselectcases](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/cmd/compile/internal/gc/select.go#L108)：

- select 不存在任何的 case；
- select 只存在一个 case；
- select 存在两个 case，其中一个 case 是 default；
- select 存在多个 case；

## 直接阻塞
当 select 结构中不包含任何 case：

```go
func walkselectcases(cases *Nodes) []*Node {
	n := cases.Len()

	if n == 0 {
		return []*Node{mkcall("block", nil, nil)}
	}
	...
}
```

这段代码直接将类似 `select {}` 的语句转换成调用 `runtime.block` 函数：

```go
func block() {
	gopark(nil, nil, waitReasonSelectNoCases, traceEvGoStop, 1)
}
```

调用 `runtime.gopark` 让出当前 Goroutine 对处理器的使用权并传入等待原因 `waitReasonSelectNoCases`。

**空的 select 语句会直接阻塞当前 Goroutine，导致 Goroutine 进入无法被唤醒的永久休眠状态**。

## 单一管道
如果当前的 select 条件只包含一个 case，那么编译器会将 select 改写成 `if` 条件语句：

```go
// 改写前
select {
case v, ok <-ch: // case ch <- v
    ...    
}

// 改写后
if ch == nil {
    block()
}
v, ok := <-ch // case ch <- v
...
```

当 case 中的 Channel 是空指针时，会直接挂起当前 Goroutine 并陷入永久休眠。

## 非阻塞操作

当 select 中仅包含两个 case，并且其中一个是 default 时，Go 语言的编译器就会认为这是一次非阻塞的收发操作。

## 常见流程 

在默认的情况下，编译器会使用如下的流程处理 select 语句：

1. 将所有的 `case` 转换成包含 Channel 以及类型等信息的 `runtime.scase` 结构体；
2. 调用运行时函数 `runtime.selectgo` 从多个准备就绪的 Channel 中选择一个可执行的 `runtime.scase` 结构体；
3. 通过 `for` 循环生成一组 `if` 语句，在语句中判断自己是不是被选中的 `case`；

```go
selv := [3]scase{}
order := [6]uint16
for i, cas := range cases {
    c := scase{}
    c.kind = ...
    c.elem = ...
    c.c = ...
}
chosen, revcOK := selectgo(selv, order, 3)
if chosen == 0 {
    ...
    break
}
if chosen == 1 {
    ...
    break
}
if chosen == 2 {
    ...
    break
}
```
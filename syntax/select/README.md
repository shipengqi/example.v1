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

Channel 的发送过程，当 case 中表达式的类型是 OSEND 时，编译器会使用条件语句和 [runtime.selectnbsend](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/chan.go#L686) 函数改写代码：

```go
select {
case ch <- i:
    ...
default:
    ...
}

if selectnbsend(ch, i) {
    ...
} else {
    ...
}
```

```go
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
	return chansend(c, elem, false, getcallerpc())
}
```

向 `runtime.chansend` 函数传入了非阻塞，所以在**不存在接收方或者缓冲区空间不足时，当前 Goroutine 都不会阻塞而是会直接返回**。

从 Channel 中接收数据可能会返回一个或者两个值，所以接收数据的情况会比发送稍显复杂

```go
// 改写前
select {
case v <- ch: // case v, ok <- ch:
    ......
default:
    ......
}

// 改写后
if selectnbrecv(&v, ch) { // if selectnbrecv2(&v, &ok, ch) {
    ...
} else {
    ...
}
```

返回值数量不同会导致使用函数的不同，两个用于非阻塞接收消息的函数 `runtime.selectnbrecv` 和 `runtime.selectnbrecv2` 只是对 `runtime.chanrecv`
返回值的处理稍有不同：

```go
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool) {
	selected, _ = chanrecv(c, elem, false)
	return
}

func selectnbrecv2(elem unsafe.Pointer, received *bool, c *hchan) (selected bool) {
	selected, *received = chanrecv(c, elem, false)
	return
}
```

与 `runtime.chansend` 一样，[runtime.chanrecv]() 也提供了一个 block 参数用于控制这次接收是否阻塞。

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

最重要的就是用于选择待执行 case 的运行时函数 [runtime.selectgo](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/runtime/select.go#L121)

执行过程：

1. 执行一些必要的初始化操作并确定 case 的处理顺序；
2. 在循环中根据 case 的类型做出不同的处理；

`runtime.selectgo` 函数首先会进行执行必要的初始化操作并决定处理 case 的两个顺序 — 轮询顺序 `pollOrder` 和加锁顺序 `lockOrder`：

```go
func selectgo(cas0 *scase, order0 *uint16, ncases int) (int, bool) {
    cas1 := (*[1 << 16]scase)(unsafe.Pointer(cas0))
    order1 := (*[1 << 17]uint16)(unsafe.Pointer(order0))

    ncases := nsends + nrecvs
    scases := cas1[:ncases:ncases]
    pollorder := order1[:ncases:ncases]
    lockorder := order1[ncases:][:ncases:ncases]

    norder := 0
    for i := range scases {
        cas := &scases[i]
    }

    for i := 1; i < ncases; i++ {
        j := fastrandn(uint32(i + 1))
        pollorder[norder] = pollorder[j]
        pollorder[j] = uint16(i)
        norder++
    }
    pollorder = pollorder[:norder]
    lockorder = lockorder[:norder]

    // 根据 Channel 的地址排序确定加锁顺序
    ...
    sellock(scases, lockorder)
    ...
}
```

runtime.selectgo 函数的主循环，它会分三个阶段查找或者等待某个 Channel 准备就绪：

查找是否已经存在准备就绪的 Channel，即可以执行收发操作；
将当前 Goroutine 加入 Channel 对应的收发队列上并等待其他 Goroutine 的唤醒；
当前 Goroutine 被唤醒之后找到满足条件的 Channel 并进行处理；
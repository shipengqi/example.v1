# for range

Go 中的经典循环在编译器由以下四个部分组成：

- 初始化循环的 Ninit；
- 循环的继续条件 Left；
- 循环体结束时执行的 Right；
- 循环体 NBody

```go
for Ninit; Left; Right {
    NBody
}
```

范围循环在 Go 语言中更常见、实现也更复杂。这种循环同时使用 `for` 和 `range` 两个关键字，编译器会在编译期间将所有 `for-range` 循环变成经典循环。

所有的 `for-range` 循环都会被 [cmd/compile/internal/gc.walkrange](https://github.com/golang/go/blob/41d8e61a6b9d8f9db912626eb2bbc535e929fefc/src/cmd/compile/internal/gc/range.go#L158) 
转换成不包含复杂结构、只包含基本表达式的语句。

## 数组和切片
对于数组和切片来说，Go 语言有三种不同的遍历方式：

- 分析遍历数组和切片清空元素的情况；
- 分析使用 `for range a {}` 遍历数组和切片，不关心索引和数据的情况；
- 分析使用 `for i := range a {}` 遍历数组和切片，只关心索引的情况；
- 分析使用 `for i, elem := range a {}` 遍历数组和切片，关心索引和数据的情况；

```go
	switch t.Etype {
	default:
		Fatalf("walkrange")

	case TARRAY, TSLICE:
		// 优化遍历数组或者切片并删除全部元素的逻辑
		if arrayClear(n, v1, v2, a) {
			lineno = lno
			return n
		}
```

```go
// 原代码
for i := range a {
	a[i] = zero
}

// 优化后
if len(a) != 0 {
	hp = &a[0]
	hn = len(a)*sizeof(elem(a))
	memclrNoHeapPointers(hp, hn)
	i = len(a) - 1
}
```

相比于依次清除数组或者切片中的数据，Go 语言会直接使用 `runtime.memclrNoHeapPointers` 或者 `runtime.memclrHasPointers` 清除目标数组
内存空间中的全部数据，并在执行完成后更新遍历数组的索引。

ORANGE 节点的处理：

```go
		ha := a

		hv1 := temp(types.Types[TINT])
		hn := temp(types.Types[TINT])

		init = append(init, nod(OAS, hv1, nil))
		init = append(init, nod(OAS, hn, nod(OLEN, ha, nil)))

		n.Left = nod(OLT, hv1, hn)
		n.Right = nod(OAS, hv1, nod(OADD, hv1, nodintconst(1)))

		if v1 == nil {
			break
		}
```

如果循环是 `for range a {}`，那么就满足了上述代码中的条件 `v1 == nil`，即循环不关心数组的索引和数据，这种循环会被编译器转换成如下形式：

```go
ha := a
hv1 := 0
hn := len(ha)
v1 := hv1
for ; hv1 < hn; hv1++ {
    ...
}
```
如果在遍历数组时需要使用索引 `for i := range a {}`，那么编译器会继续会执行下面的代码：

```go
		if v2 == nil {
			body = []*Node{nod(OAS, v1, hv1)}
			break
		}
```

`v2 == nil` 意味着调用方不关心数组的元素，只关心遍历数组使用的索引。与第一种循环相比，这种循环在循环体中添加了 `v1 := hv1` 语句，传
递遍历数组时的索引：

```go
ha := a
hv1 := 0
hn := len(ha)
v1 := hv1
for ; hv1 < hn; hv1++ {
    v1 = hv1
    ...
}
```

同时去遍历索引和元素也很常见。处理这种情况会使用下面这段的代码：

```go
		tmp := nod(OINDEX, ha, hv1)
		tmp.SetBounded(true)
		a := nod(OAS2, nil, nil)
		a.List.Set2(v1, v2)
		a.Rlist.Set2(hv1, tmp)
		body = []*Node{a}
	}
	n.Ninit.Append(init...)
	n.Nbody.Prepend(body...)

	return n
```

使用者同时关心索引和切片的情况。它不仅会在循环体中插入更新索引的语句，还会插入赋值操作让循环体内部的代码能够访问数组中的元素：
```go
ha := a
hv1 := 0
hn := len(ha)
v1 := hv1
v2 := nil
for ; hv1 < hn; hv1++ {
    tmp := ha[hv1]
    v1, v2 = hv1, tmp
    ...
}
```

对于所有的 range 循环，Go 语言都会在编译期将原切片或者数组赋值给一个新变量 `ha`，在赋值的过程中就发生了**拷贝**，而我们又通过 `len` 关键字预先
获取了切片的长度，所以**在循环中追加新的元素也不会改变循环执行的次数**。

**循环中使用的这个变量 `v2` 会在每一次迭代被重新赋值而覆盖，赋值时也会触发拷贝**。
```go
	arr := []int{1, 2, 3}
	var newArr []*int
	for _, v := range arr {
		newArr = append(newArr, &v)
	}
	for _, v := range newArr {
		fmt.Println(*v)
	}
```

因为在循环中获取返回变量的地址都完全相同，所以输出会是 `3 3 3`。

## Map

在遍历 map 时，编译器会使用 `runtime.mapiterinit` 和 `runtime.mapiternext` 两个运行时函数重写原始的 `for-range` 循环：

```go
ha := a
hit := hiter(n.Type)
th := hit.Type
mapiterinit(typename(t), ha, &hit)
for ; hit.key != nil; mapiternext(&hit) {
    key := *hit.key
    val := *hit.val
}
```

上述代码是展开 `for key, val := range hash {}` 后的结果。

三种不同的情况分别向循环体插入了不同的赋值语句。
- `for range hash {}`
- `for k := range hash {}`
- `for k, v := range hash {}` 

遍历哈希表时会使用 `runtime.mapiterinit` 函数**初始化遍历开始的元素**：
```go
func mapiterinit(t *maptype, h *hmap, it *hiter) {
	it.t = t
	it.h = h
	it.B = h.B
	it.buckets = h.buckets

	r := uintptr(fastrand())
	it.startBucket = r & bucketMask(h.B)
	it.offset = uint8(r >> h.B & (bucketCnt - 1))
	it.bucket = it.startBucket
	mapiternext(it)
}
```

初始化 `runtime.hiter` 结构体中的字段，并通过 **`runtime.fastrand` 生成一个随机数帮助我们随机选择一个遍历桶的起始位置**。
Go 团队在设计哈希表的遍历时就不想让使用者依赖固定的遍历顺序，所以引入了随机数保证遍历的随机性。

## 字符串

遍历字符串的过程与数组、切片和哈希表非常相似，只是在遍历时会获取字符串中索引对应的字节并将字节转换成 `rune`。我们在遍历字符串时拿到的值都是 `rune` 类
型的变量，`for i, r := range s {}` 的结构都会被转换成如下所示的形式：

```go
ha := s
for hv1 := 0; hv1 < len(ha); {
    hv1t := hv1
    hv2 := rune(ha[hv1])
    if hv2 < utf8.RuneSelf {
        hv1++
    } else {
        hv2, hv1 = decoderune(ha, hv1)
    }
    v1, v2 = hv1t, hv2
}
```

## chan

`for v := range ch {}` 的语句最终会被转换成如下的格式：
```go
ha := a
hv1, hb := <-ha
for ; hb != false; hv1, hb = <-ha {
    v1 := hv1
    hv1 = nil
    ...
}
```

该循环会使用 `<-ch` 从管道中取出等待处理的值，这个操作会调用 `runtime.chanrecv2` 并阻塞当前的协程，当 `runtime.chanrecv2` 返回时会
根据布尔值 `hb` 判断当前的值是否存在：

- 如果不存在当前值，意味着当前的管道已经被关闭；
- 如果存在当前值，会为 `v1` 赋值并清除 `hv1` 变量中的数据，然后重新陷入阻塞等待新数据；

`hb != false` 也就意味着，终止条件是 chan 关闭，这就是为什么 range 可以持续从一个通道接收值，直到 chan 关闭。如果 chan 没有元素，
或者 chan 为 nil，则会阻塞 goroutine
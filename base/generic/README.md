# 泛型

Go 泛型泛型语法为中括号 `[]`：

```go
func Print[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}
```

**类型参数列表**出现在常规参数之前。类型参数列表使用中括号 `[]`。

```
Print [int]([] int{1, 2, 3})
// Output
// 1
// 2
// 3
```

Java、C++、Rust 等大多数语言中，使用 `F<T>` 来标识泛型，Go 为什么使用中括号？

```go
a, b := w < x, y > (z)
```

上面的代码中，如果没有类型信息，则不能区分右值是一对表达式 `w < x , y > (z)`，还是返回两个结果值的函数调用，而在此种情况下，Go 希望在没有类型信息的情况
下也能进行正确的解析，那尖括号则满足不了这个条件。


## 类型约束

类型约束用于进一步限制可以作为 T 的可能类型。Go 本身提供了一些预定义的类型约束，但也可以使用自定义的类型约束。

```go
type Node[T any] struct {
    value T
}
```

任意类型（`any`）约束允许 T 实际上是任何类型（`any` 实际上就是 `interface{}` 的别名）。如果节点值需要进行比较，有一个 `comparable` 类型约束，满足这个预定义约束的类型可以使用 `==` 进行比较。

```go
type Node[T comparable] struct {
    value T
}
```

`comparable` 比较容易引起误解的一点是很多人容易把他与可排序搞混淆。可比较指的是可以执行 `!=` `==` 操作的类型，并没确保这个类型可以执行大小比较（ `>,<,<=,>=` ）。
Go 语言并没有像 `comparable` 这样直接内置对应的关键词，所以想要的话需要自己来定义相关接口，可以参考 `golang.org/x/exp/constraints` 的定义：

```go
// Ordered 代表所有可比大小排序的类型
type Ordered interface {
	Integer | Float | ~string
}
```

**近似约束元素**：

碰到下面这样的类型别名：

```go
type Phone string
type Email string
type Address string
```

Go 1.18 中扩展了近似约束元素，以上述例子来说，即：基础类型为 `string` 的类型。语法表现为：

```go
type String interface{~string}
```

`String` 类型即可表示 `Phone | Email | Address`，使用约束：

```go
func Deduplicate[T String] (str T) string{
   // do something
}
```

**任何类型都可以作为一个类型约束**。Go 1.18 引入了一种新的 `interface` 语法，可以嵌入其他数据类型。

```go
type Numeric interface {
    int | float32 | float64
}
```

这意味着一个接口不仅可以定义一组方法，还可以定义一组类型。使用 `Numeric` 接口作为类型约束，意味着值可以是整数或浮点数。

```go
type Node[T Numeric] struct {
    value T
}
```

只能使用确定的类型进行联合类型约束：

```go
 // GOOD
func PrintInt64OrFloat64[T int64|float64](t T) {
   fmt.Printf( "%v\n" , t)
}

type someStruct struct {}

// GOOD
func PrintInt64OrSomeStruct[T int64|*someStruct](t T) {
   fmt.Printf( "t: %v\n" , t)
}

// BAD，不能在联合类型中使用 ，且不能通过编译
func handle[T io.Closer | Flusher](t T) {
   err := t.Flush()
   if err != nil {
      fmt.Println( "failed to flush: " , err.Error())
   }

   err = t.Close()
   if err != nil {
      fmt.Println( "failed to close: " , err.Error())
   }
}

type Flusher interface {
   Flush() error
}
```


Go 已经有一个接近于我们需要的约束的构造：接口类型。接口类型是一组方法。

类型约束定义：
```go
type Stringer interface {
   String() string
}
```

约束使用：
```go
func Stringify[T Stringer](s []T) (ret []string) {
   for _, v := range s {
      ret = append(ret, v.String())
   }
   return ret
}
```

支持多个类型参数和约束：
```go
type Stringer interface {
   String() string
}

type Plusser interface {
   Plus(string) string
}

func ConcatTo[S Stringer, P Plusser](s []S, p []P) []string {
   r := make([]string, len(s))
   for i, v := range s {
      r[i] = p[i].Plus(v.String())
   }
   return r
}
```

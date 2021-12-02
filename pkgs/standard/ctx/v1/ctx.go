package main

import (
	"context"
	"fmt"
	"time"
)

// context 包中包含了四个用于繁衍 Context 值的函数。
// WithValue
// WithCancel
// WithDeadline
// WithTimeout
// 其中的 WithCancel、WithDeadline 和 WithTimeout 都是被用来基于给定的 Context 值产生可撤销的子值的。
//
// WithCancel 在被调用后会产生两个结果值。第一个结果值就是那个可撤销的 Context 值，而第二个结果值则是用于触发撤销信号的函数。
//            在撤销函数被调用之后，对应的 Context 值会先关闭它内部的接收通道，也就是它的 Done 方法会返回的那个通道。
//            然后，它会向它的所有子值（或者说子节点）传达撤销信号。这些子值会如法炮制，把撤销信号继续传播下去。
//            最后，这个 Context 值会断开它与其父值之间的关联。
//
// WithDeadline 和 WithTimeout 生成的 Context 值也是可撤销的。它们不但可以被手动撤销，还会依据在生成时被给定的过期时间，自动地进行定时撤销。
// 这里定时撤销的功能是借助它们内部的计时器来实现的。
// 当过期时间到达时，这两种 Context 值的行为与 Context 值被手动撤销时的行为是几乎一致的，只不过前者会在最后停止并释放掉其内部的计时器。
//
// 注意，通过调用 context.WithValue 函数得到的 Context 值是不可撤销的。撤销信号在被传播时，若遇到它们则会直接跨过，并试图将信号直接传给它们的子值。
//
// WithTimeout 很常见的另外一个方法，是便捷操作。实际上是对于 WithDeadline 的封装
//
// func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
//	    return WithDeadline(parent, time.Now().Add(timeout))
// }
//
// WithValue 函数在产生新的 Context 值（以下简称含数据的Context值）的时候需要三个参数，即：父值、键和值。与“字典对于键的约束”类似，这里键的类型必须是可判等的。
// 不过，这种 Context 值并不是用字典来存储键和值的，后两者只是被简单地存储在前者的相应字段中而已。
//
// Context 类型的 Value 方法就是被用来获取数据的。在我们调用含数据的 Context 值的 Value 方法时，它会先判断给定的键，是否与当前值中存储的
// 键相等，如果相等就把该值中存储的值直接返回，否则就到其父值中继续查找。如果其父值中仍然未存储相等的键，那么该方法就会沿着上下文根节点的方向一路查找下去。
// 其他几种 Context 值都是无法携带数据的。因此，Context 值的 Value 方法在沿路查找的时候，会直接跨过那几种值。
//
// Context 接口并没有提供改变数据的方法。因此，在通常情况下，我们只能通过在上下文树中添加含数据的 Context 值来存储新的数据，或者通过撤销此种值的父
// 值丢弃掉相应的数据。如果你存储在这里的数据可以从外部改变，那么必须自行保证安全。

// Context 接口
// type Context interface {
//     Deadline() (deadline time.Time, ok bool)
//
//     Done() <-chan struct{}
//
//     Err() error
//
//     Value(key interface{}) interface{}
// }
// Deadline 方法是获取设置的截止时间的意思，第一个返回值截止时间，到了这个时间点，Context 会自动发起取消请求；
//          第二个返回值 ok==false 时表示没有设置截止时间，如果需要取消的话，需要调用取消函数进行取消。
//
// Done 方法返回一个只读的 chan，类型为 struct{}，我们在 goroutine 中，如果该方法返回的 chan 可以读取，则意味着 parent context 已经发
// 起了取消请求，我们通过 Done 方法收到这个信号后，就应该做清理操作，然后退出 goroutine，释放资源。
//
// Err 方法返回取消的错误原因，因为什么 Context 被取消。
//
// Value 方法获取该 Context 上绑定的值，是一个键值对，所以要通过一个Key才可以获取对应的值，这个值一般是线程安全的。

// Context 使用原则：
// 不要把 Context 放在结构体中，要以参数的方式传递
// 以 Context 作为参数的函数方法，应该把 Context 作为第一个参数，放在第一位。
// 给一个函数方法传递 Context 的时候，不要传递 nil，如果不知道传递什么，就使用 `context.TODO`
// Context 的 Value 相关方法应该传递必须的数据，不要什么数据都使用这个传递
// Context 是线程安全的，可以放心的在多个 goroutine 中传递

type myKey int

func main() {
	keys := []myKey{
		myKey(20),
		myKey(30),
		myKey(60),
		myKey(61),
	}
	values := []string{
		"value in node2",
		"value in node3",
		"value in node6",
		"value in node6Branch",
	}

	rootNode := context.Background()
	node1, cancelFunc1 := context.WithCancel(rootNode)
	defer cancelFunc1()

	// 示例 1
	node2 := context.WithValue(node1, keys[0], values[0])
	node3 := context.WithValue(node2, keys[1], values[1])
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[0], node3.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[1], node3.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node3: %v\n",
		keys[2], node3.Value(keys[2]))
	fmt.Println()

	// 示例 2
	node4, _ := context.WithCancel(node3)
	node5, _ := context.WithTimeout(node4, time.Hour)
	fmt.Printf("The value of the key %v found in the node5: %v\n",
		keys[0], node5.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node5: %v\n",
		keys[1], node5.Value(keys[1]))
	fmt.Println()

	// 示例 3
	node6 := context.WithValue(node5, keys[2], values[2])
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		keys[0], node6.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node6: %v\n",
		keys[2], node6.Value(keys[2]))
	fmt.Println()

	// 示例 4
	node6Branch := context.WithValue(node5, keys[3], values[3])
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[1], node6Branch.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[2], node6Branch.Value(keys[2]))
	fmt.Printf("The value of the key %v found in the node6Branch: %v\n",
		keys[3], node6Branch.Value(keys[3]))
	fmt.Println()

	// 示例 5
	node7, _ := context.WithCancel(node6)
	node8, _ := context.WithTimeout(node7, time.Hour)
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[1], node8.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[2], node8.Value(keys[2]))
	fmt.Printf("The value of the key %v found in the node8: %v\n",
		keys[3], node8.Value(keys[3]))

	// 示例 6
	node7Branch := context.WithValue(node8, keys[3], "overwrite")
	fmt.Printf("The value of the key %v found in the node7Branch: %v\n",
		keys[1], node7Branch.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node7Branch: %v\n",
		keys[2], node7Branch.Value(keys[2]))
	fmt.Printf("The value of the key %v found in the node7Branch: %v\n",
		keys[3], node7Branch.Value(keys[3]))
	fmt.Println()
}

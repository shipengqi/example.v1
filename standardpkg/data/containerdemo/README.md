# container

container 包包现了三个复杂的数据结构：堆，链表，环。 就意味着使用这三个数据结构的时候不需要再自己实现。

##  链表

链表就是一个有 `prev` 和 `next` 指针的数组了。

`container` 包中有两个公开的结构—— `List` 和 `Element`，`List` 实现了一个双向链表（简称链表）， 而 `Element` 则代表了链表中元素的结构。

```go
type Element struct {
	next, prev *Element  // 上一个元素和下一个元素
	list *List  // 元素所在链表
	Value interface{}  // 元素
}

type List struct {
	root Element  // 链表的根元素
	len  int      // 链表的长度
}
```

`List` 的方法:

```go
type Element
    func (e *Element) Next() *Element
    func (e *Element) Prev() *Element
type List
    func New() *List
    func (l *List) Back() *Element   // 最后一个元素
    func (l *List) Front() *Element  // 第一个元素
    func (l *List) Init() *List  // 链表初始化 清除 list
    func (l *List) InsertAfter(v interface{}, mark *Element) *Element // 在某个元素后插入
    func (l *List) InsertBefore(v interface{}, mark *Element) *Element  // 在某个元素前插入
    func (l *List) Len() int // 链表长度
    func (l *List) MoveAfter(e, mark *Element)  // 把 e 元素移动到 mark 元素之后
    func (l *List) MoveBefore(e, mark *Element)  // 把 e 元素移动到 mark 元素之前
    func (l *List) MoveToBack(e *Element) // 把 e 元素移动到队列最后
    func (l *List) MoveToFront(e *Element) // 把 e 元素移动到队列最头部
    func (l *List) PushBack(v interface{}) *Element  // 在队列最后插入元素
    func (l *List) PushBackList(other *List)  // 在队列最后插入接上新队列
    func (l *List) PushFront(v interface{}) *Element  // 在队列头部插入元素
    func (l *List) PushFrontList(other *List) // 在队列头部插入接上新队列
    func (l *List) Remove(e *Element) interface{} // 删除某个元素
```

##  堆

### 什么是堆
堆（Heap，也叫优先队列）是计算机科学中一类特殊的数据结构的统称。**堆通常是一个可以被看做一棵树的数组对象**。

堆具有以下特性：
- 任意节点小于（或大于）它的所有后裔，最小元（或最大元）在堆的根上（堆序性）。
- 堆总是一棵完全树。即除了最底层，其他层的节点都被元素填满，且最底层尽可能地从左到右填入。

将根节点最大的堆叫做最大堆或大根堆，根节点最小的堆叫做最小堆或小根堆。

### heap 包提供的方法

1. `func Init(h Interface)` 对 heap 进行初始化（堆化），生成小根堆（或大根堆）
2. `func Pop(h Interface) interface{}` 往堆里面插入元素
3. `func Push(h Interface, x interface{})` 从堆顶 pop 出元素
4. `func Remove(h Interface, i int) interface{}` 从指定位置删除数据，并返回删除的数据
5. `func Fix(h Interface, i int)` 从 i 位置数据发生改变后，对堆再平衡，优先级队列使用到了该方法
5. `type Interface`

`Interface` 接口：

```go
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

除了堆接口定义的两个方法：

```go
Push(x interface{})
Pop() interface{}
```

还继承了 `sort.Interface`, 需要实现三个方法：

```go
Len() int
Less(i, j int) bool
Swap(i, j int)
```

**通过对 `heap.Interface` 中的 `Less` 方法的不同实现，来实现最大堆和最小堆**。

下面是最小堆的示例：
```go
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	for h.Len() > 0 {
		fmt.Printf("%d \n", heap.Pop(h))
    }
}
```

如果要实现最大堆，只需要修改 `Less` 方法的实现：

```go
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
```

## 环

`Ring` 类型代表环形链表的一个元素，同时也代表链表本身。

```go
type Ring struct {
	next, prev *Ring
	Value      interface{}
}
```

初始化环的时候，需要定义好环的大小，然后对环的每个元素进行赋值。环还提供一个 `Do` 方法，能遍历一遍环，对每个元素执行
一个 `function`。

ring 提供的方法有

```go
type Ring
    func New(n int) *Ring // 创建一个长度为 n 的环形链表
    func (r *Ring) Do(f func(interface{})) // 遍历环形链表中的每一个元素 x 进行 f(x) 操作
    func (r *Ring) Len() int // 获取环形链表长度
    
    // Link 连接 r 和 s，并返回 r 原本的后继元素 r.Next()。r 不能为空。
    // 如果 r 和 s 在同一环形链表中，则删除 r 和 s 之间的元素，
    // 被删除的元素组成一个新的环形链表，返回值为该环形链表的指针（即删除前，r->Next() 表示的元素）
    // 如果 r 和 s 不在同一个环形链表中，则将 s 插入到 r 后面，返回值为插入 s 后，
    // s 最后一个元素的下一个元素（即插入前，r->Next() 表示的元素）
    func (r *Ring) Link(s *Ring) *Ring

    func (r *Ring) Move(n int) *Ring // 返回移动 n 个位置（n>=0 向前移动，n<0 向后移动）后的元素
    func (r *Ring) Next() *Ring // 返回下一个元素
    func (r *Ring) Prev() *Ring // 返回前一个元素
    func (r *Ring) Unlink(n int) *Ring // 删除 r 后面的 n % r.Len() 个元素，如果 n % r.Len() == 0，不修改 r。
```

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

`container` 包的堆使用的数据结构是最小二叉树。即根节点比左边子树和右边子树的所有值都小。go 的堆包只是实现了一个接口：

```go
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

这个堆结构继承自 sort.Interface, 回顾下 sort.Interface，它需要实现三个方法


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

示例：
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

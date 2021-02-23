package main

import (
	"container/heap"
	"fmt"
)

type Item struct {
	value    string // 优先级队列中的数据，可以是任意类型，这里使用 string
	priority int    // 优先级队列中节点的优先级
	index    int    // index 是该节点在堆中的位置
}

// 优先级队列需要实现 heap 的 Interface
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

// 这里用的是小于号，生成的是最小堆
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// 将 index 置为 -1 是为了标识该数据已经出了优先级队列了
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	item.index = -1
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// 更新修改了优先级和值的 item 在优先级队列中的位置
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

func main() {
	// 创建节点并设计他们的优先级
	items := map[string]int{"二毛": 5, "张三": 3, "狗蛋": 9}
	i := 0
	pq := make(PriorityQueue, len(items)) // 创建优先级队列，并初始化
	for k, v := range items {             // 将节点放到优先级队列中
		pq[i] = &Item{
			value:    k,
			priority: v,
			index:    i,
		}
		i++
	}

	heap.Init(&pq) // 初始化堆
	item := &Item{ // 创建一个 item
		value:    "李四",
		priority: 1,
	}
	heap.Push(&pq, item)           // 入优先级队列
	pq.update(item, item.value, 6) // 更新 item 的优先级
	for len(pq) > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s index:%.2d\n", item.priority, item.value, item.index)
	}
}


// Output：
// 03:张三 index:-01
// 05:二毛 index:-01
// 06:李四 index:-01
// 09:狗蛋 index:-01

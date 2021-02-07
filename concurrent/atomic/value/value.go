package main

import (
	"fmt"
	"sync/atomic"
)

// 原子值的两个限制：
// 1. 作为参数传入到 Store 的值不能为 nil
// 2. 传入 Store 的值要与之前传入的值类型相同
// 注意：atomic.Value 不应该被复制到别的地方，比如：
// 1. 作为源值赋值给别的变量
// 2. 作为参数传入函数
// 3. 作为函数返回值
// 4. 通过 chan 传递
// 以上都会造成值的复制。可以使用指针来避免这种错误。
// 根本原因：**对结构体的复制会生成值的副本，还会生成字段的副本，如此本应该施加于此的并发安全保护就失效了**。
// 并且向副本值的操作也与原值无关
func main() {
	var v atomic.Value
	v.Store([]int{1, 3, 5, 7})
	otherStore(v)
	fmt.Printf("count value: %d\n", v.Load())
	otherStore2(&v)
	fmt.Printf("count value pointer: %d\n", v.Load())
	v2 := v
	v2.Store([]int{1, 3, 5, 7})
	fmt.Printf("count value variable assignment : %d\n", v.Load())
	ts := []int{1, 2, 3, 4}
	sliceTest(ts)
	fmt.Println("slice result", ts)
	tm := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	mapTest(tm)
	fmt.Println("map result", tm)
}

func otherStore(v atomic.Value) {
	v.Store([]int{2, 4, 6, 8})
}

// 输出 count value: [1 3 5 7]
// 使用 go vet 检查此类复制问题，为了避免此类问题，可以传递指针，如 otherStore2

func otherStore2(v *atomic.Value) {
	v.Store([]int{2, 4, 6, 8})
}

// 下面的函数只是为了测试指针传递
func mapTest(m map[string]string) {
	m["key1"] = "value100"
	fmt.Println("map index", m)
	m = map[string]string{
		"key101": "value101",
		"key102": "value102",
	}
	fmt.Println("map variable assignment", m)
}

func sliceTest(v []int) {
	v[0] = 100
	fmt.Println("slice index", v)
	v = append(v, 101)
	fmt.Println("slice append", v)
	v = []int{102, 103}
	fmt.Println("slice variable assignment", v)
}



package main

import (
	"fmt"
	"sort"
)

func main() {
	x := 11
	s := []int{3, 6, 8, 11, 45} // 已经升序排序的集合
	pos := sort.Search(len(s), func(i int) bool { return s[i] >= x })
	if pos < len(s) && s[pos] == x {
		fmt.Println(x, " 在 s 中的位置为：", pos)
	} else {
		fmt.Println("s 不包含元素 ", x)
	}
}



// Output:
// Default:
//	 [{alan 95} {hikerell 91} {acmfly 96} {leao 90}]
// IS Sorted?
//	 true
// Sorted:
//	 [{acmfly 96} {alan 95} {hikerell 91} {leao 90}]
// sort.Reverse:
//	 &{[{acmfly 96} {alan 95} {hikerell 91} {leao 90}]}
package main

import (
	"fmt"
	"sort"
)

type StuScore struct {
	name  string // 姓名
	score int    // 成绩
}

type StuScores []StuScore

func (s StuScores) Len() int {
	return len(s)
}

func (s StuScores) Less(i, j int) bool {
	return s[i].score < s[j].score
}

func (s StuScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	stus := StuScores{
		{"alan", 95},
		{"hikerell", 91},
		{"acmfly", 96},
		{"leao", 90},
	}

	// 未排序的 stus 数据
	fmt.Println("Default:\n\t", stus)
	// 调用 Sort 函数进行排序
	sort.Sort(stus)
	// 判断是否已经排好顺序
	fmt.Println("IS Sorted?\n\t", sort.IsSorted(stus))
	// 排序后的 stus 数据
	fmt.Println("Sorted:\n\t", stus)
}

// Output:
// Default:
//	 [{alan 95} {hikerell 91} {acmfly 96} {leao 90}]
// IS Sorted?
//	 true
// Sorted:
//	 [{leao 90} {hikerell 91} {alan 95} {acmfly 96}]
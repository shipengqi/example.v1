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

// v1 的 Less 实现的是升序排序，如果要得到降序排序结果，其实只要修改 Less() 函数：
func (s StuScores) Less(i, j int) bool {
	return s[i].score > s[j].score
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

	reversed := sort.Reverse(stus)
	// 需要再次调用 Sort，才会颠倒集合的排序，因为 Reverse 函数其实只是修改了 Less 函数
	sort.Sort(reversed)

	// reverse 后的 stus 数据
	fmt.Println("sort.Reverse:\n\t", reversed)

	x := 95
	searched := sort.Search(len(stus), func(i int) bool {
		return stus[i].score >= x
	})

	fmt.Printf("name: %s, score is %d\n\t", stus[searched].name, stus[searched].score)
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
// name: alan, score is 95
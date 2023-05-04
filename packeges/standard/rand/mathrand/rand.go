package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	num := 10
	for j := 0; j < num; j++ {
		res := getRand(num)
		fmt.Println(res)
	}

	fmt.Println("get rand with seed")
	for j := 0; j < num; j++ {
		res := getRandWithSeed(num)
		fmt.Println(res)
		// 为了避免运行太快导致 time.Now().UnixNano() 获取的值一样
		// 调用 rand.Seed(x) 时，x 一样，会导致每次的 rand.Intn() 也是一样的。
		// 生产环境下，强烈不推荐这种，因为高并发的情况下纳秒也可能重复。
		time.Sleep(time.Nanosecond*2)
	}
}

func getRand(num int) int {
	v := rand.Intn(num)
	return v
}

func getRandWithSeed(num int) int {
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(num)
	return v
}

// getRand 每次输出可能都是下面的数字， 这是因为没有设定随机种子：
//1
//7
//7
//9
//1
//8
//5
//0
//6
//0
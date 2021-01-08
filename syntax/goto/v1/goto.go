package main

import "fmt"

// goto 句可以无条件地转移到程序中指定的行。
// goto 语句通常与条件语句配合使用。可用来实现条件转移，跳出循环体等功能。
// 在 Go 中一般不推荐使用 goto 语句， 以免造成程序流程的混乱，使理解和调试程序困难。

func main() {
	var intervalTimes int
Loop:
	fmt.Println("test")
	for a := 0; a < 5; a++ {
		intervalTimes ++
		if intervalTimes > 10 {
			break
		}
		fmt.Println(a)
		if a == 3 {
			goto Loop
		}
	}
}

// Output：
//test
//0
//1
//2
//3
//test
//0
//1
//2
//3
//test
//0
//1
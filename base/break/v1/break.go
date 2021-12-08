package main

import "fmt"

// break 和 continue 都支持 label，使用方式也是一样的
// 在没有使用 Loop 标签的时候 break 只是跳出了第一层 for 循环
// 使用标签后跳出到指定的标签,break 只能跳出到之前，如果将 Loop 标签放在后边则会报错
// break 标签只能用于 for 循环，跳出后不再执行标签对应的 for 循环

func main() {
Loop:
	for j := 0; j < 3; j++ {
		fmt.Println("loop1", j)
		for a := 0; a < 5; a++ {
			fmt.Println("loop2", a)
			if a > 3 {
				break Loop
			}
		}
	}
}

// Output:
// loop1 0
// loop2 0
// loop2 1
// loop2 2
// loop2 3
// loop2 4

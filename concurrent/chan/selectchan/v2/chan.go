package main

import "fmt"

// 如果只有一个 case 满足，runtime 就会运行 case 中对应的语句
// 如果同时满足多个 case，runtime 会通过一个伪随机算法选中一个 case
func main() {
	chanCap := 5
	intChan := make(chan int, chanCap)
	for i := 0; i < chanCap; i++ {
		select {
		case intChan <- 1:
		case intChan <- 2:
		case intChan <- 3:
		}
	}
	for i := 0; i < chanCap; i++ {
		fmt.Printf("get form chan [%d]\n", <-intChan)
	}
}

// Output:
//get form chan [3]
//get form chan [2]
//get form chan [1]
//get form chan [3]
//get form chan [1]
// 上面的输出是不确定的，比如再次运行，输出：
//get form chan [1]
//get form chan [1]
//get form chan [1]
//get form chan [3]
//get form chan [3]

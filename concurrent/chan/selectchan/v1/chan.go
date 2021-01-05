package main

import (
	"fmt"
)

// select 和 switch 语句看起来非常类似，select 只能用于 chan 发送和接收操作
// case 后面只能是 chan 的发送或者接收语句
// select case 后面的表达式会先进行求值，求值顺序是从左到右，自上而下
// default 语句非常重要，如果没有满足的 case，也没有 default，当前 goroutine 就会阻塞
var intChan1 chan int
var intChan2 chan int
var channels = []chan int{intChan1, intChan2}
var numbers = []int{1, 2, 3, 4, 5}

func main() {
	select {
	// intChan1 和 intChan2 都未初始化，所以是 nil，发送操作都会被阻塞，所以最后进入了 default 分支
	case getChan(0) <- getNum(0):
		fmt.Println("case1")
	case getChan(1) <- getNum(1):
	    fmt.Println("case2")
	default:
		fmt.Println("default")
	}
}

func getNum(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]
}
func getChan(i int) chan int {
	fmt.Printf("channels[%d]\n", i)
	return channels[i]
}
// Output:
//channels[0]
//numbers[0]
//channels[1]
//numbers[1]
//default
// 上面四行输出说明了 case 后面的表达式进行求值计算的顺序。


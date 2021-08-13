package main

import "fmt"

func main() {
	{
		defer fmt.Println("defer runs")
		fmt.Println("block ends")
	}

	fmt.Println("main ends")
}

// Output:
//block ends
//main ends
//defer runs
// defer 传入的函数不是在退出代码块的作用域时执行的，它只会在当前函数和方法返回之前被调用。
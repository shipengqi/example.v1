package main

// select 语句在执行时会遇到以下两种情况：
// 当存在可以收发的 Channel 时，直接处理该 Channel 对应的 case；
// 当不存在可以收发的 Channel 时，执行 default 中的语句；

func main()  {
	ch := make(chan int)
	select {
	case i := <-ch:
		println(i)

	default:
		println("default")
	}
}

// Output:
// default
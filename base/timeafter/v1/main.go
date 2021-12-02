package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	ch := make(chan int, 10)

	go func() {
		var i = 1
		for {
			i++
			ch <- i
		}
	}()

	go func() {
		// 开启 pprof，监听请求
		ip := "127.0.0.1:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
		}
	}()

	var count int
	for {
		select {
		case x := <-ch:
			fmt.Println("ch: ", x)
		case <-time.After(3 * time.Minute):
			fmt.Println("timeout: ", time.Now().Unix())
		}
		count++
		fmt.Println("count: ", count)
	}
}

// 此段代码会引发内存泄漏
// time.After 会调用 time.NeWTimer 的不断创建新的 timer 并内存申请
// 因为 for在循环时，就会调用都 select 语句，因此在每次进行 select 时，都会重新初始化一个全新的计时器（Timer）。
// 这个计时器在 3 分钟后，才会被激活，但是激活后已经跟 select 无引用关系（因为 case x := <-ch 接收到了值），因此很合理的也就被 GC 给清理掉了。
// 但是 被抛弃的 time.After 的定时任务还是在时间堆中等待触发，在定时任务未到期之前，是不会被 GC 清除的。
// 也就是说每次循环实例化的新定时器对象需要 3 分钟才会可能被 GC 清理掉，如果我们把上面定时器时间改成 10 秒钟，会有所改善。

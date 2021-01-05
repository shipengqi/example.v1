package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)

func main() {
	c := make(chan struct{}, 2)
	go func() {
		for {
			if elem, ok := <- mapChan; ok {
				elem["count"] ++
			} else {
				break
			}
		}
		fmt.Println("stopped. [receiver]")
		c <- struct{}{}
	}()

	go func() {
		countMap := make(map[string]int)
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("count map: %+v. [sender]\n", countMap)
		}
		close(mapChan)
		c <- struct{}{}
	}()

	<- c
	<- c
}

// Output:
//count map: map[count:1]. [sender]
//count map: map[count:2]. [sender]
//count map: map[count:3]. [sender]
//count map: map[count:4]. [sender]
//count map: map[count:5]. [sender]
//stopped. [receiver]
// 上面的示例中 countMap 的值被另一个 goroutine 中的操作改变了，也就是说接收方对 chan 中的副本修改，影响了原值。
// 这是因为 map 是引用类型，chan 中的副本只是指针的拷贝。如果是值类型，就不会出现这种情况。

package main

import (
	"fmt"
	"time"
)

type Counter struct {
	count int
}

func (c *Counter) String() string {
	return fmt.Sprintf("{count: %d}", c.count)
}

// var mapChan = make(chan map[string]Counter, 1)
var mapChan = make(chan map[string]*Counter, 1)

func main() {
	c := make(chan struct{}, 2)
	go func() {
		for {
			if elem, ok := <- mapChan; ok {
				 counter := elem["count"]
				 counter.count ++
			} else {
				break
			}
		}
		fmt.Println("stopped. [receiver]")
		c <- struct{}{}
	}()

	go func() {
		//countMap := map[string]Counter{
		//	"count": Counter{},
		//}
		countMap := map[string]*Counter{
			"count": &Counter{},
		}
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
//count map: map[count:{count:0}]. [sender]
//count map: map[count:{count:0}]. [sender]
//count map: map[count:{count:0}]. [sender]
//count map: map[count:{count:0}]. [sender]
//count map: map[count:{count:0}]. [sender]
//stopped. [receiver]

// 如果把 map 的 Counter 改成 *Counter，Output:
//count map: map[count:0xc0000a0068]. [sender]
//count map: map[count:0xc0000a0068]. [sender]
//count map: map[count:0xc0000a0068]. [sender]
//count map: map[count:0xc0000a0068]. [sender]
//count map: map[count:0xc0000a0068]. [sender]
//stopped. [receiver]
// 实现 String 方法后，Output:
//count map: map[count:{count: 1}]. [sender]
//count map: map[count:{count: 2}]. [sender]
//count map: map[count:{count: 3}]. [sender]
//count map: map[count:{count: 4}]. [sender]
//count map: map[count:{count: 5}]. [sender]
//stopped. [receiver]

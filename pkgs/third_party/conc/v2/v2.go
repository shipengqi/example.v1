package v2

import "sync"

func process(stream chan int) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for elem := range stream {
				handle(elem)
			}
		}()
	}
	wg.Wait()
}

func handle(e int) {

}

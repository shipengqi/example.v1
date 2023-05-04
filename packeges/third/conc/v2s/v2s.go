package v2s

import "github.com/sourcegraph/conc/pool"

func process(stream chan int) {
	p := pool.New().WithMaxGoroutines(10)
	for elem := range stream {
		e := elem
		p.Go(func() {
			handle(e)
		})
	}
	p.Wait()
}

func handle(e int) {

}

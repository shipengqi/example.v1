package v2

import (
	"sync"

	"go.uber.org/atomic"
)

func concMap(
	input []int,
	f func(int) int,
) []int {
	res := make([]int, len(input))
	var idx atomic.Int64

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				i := int(idx.Add(1) - 1)
				if i >= len(input) {
					return
				}

				res[i] = f(input[i])
			}
		}()
	}
	wg.Wait()
	return res
}

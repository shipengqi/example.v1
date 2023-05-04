package v2

import (
	"github.com/sourcegraph/conc/stream"
)

func mapStream(
	in chan int,
	out chan int,
	f func(int) int,
) {
	s := stream.New().WithMaxGoroutines(10)
	for elem := range in {
		elem := elem
		s.Go(func() stream.Callback {
			res := f(elem)
			return func() { out <- res }
		})
	}
	s.Wait()
}

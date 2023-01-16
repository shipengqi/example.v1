package v2

import (
	"github.com/sourcegraph/conc/iter"
)

func concMap(
	input []int,
	f func(*int) int,
) []int {
	return iter.Map(input, f)
}

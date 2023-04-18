package v2

import (
	"github.com/sourcegraph/conc/iter"
)

// Map 与 MapErr 也只是对 ForEachIdx 的封装，区别是处理 error
func concMap(
	input []int,
	f func(*int) int,
) []int {
	return iter.Map(input, f)
}

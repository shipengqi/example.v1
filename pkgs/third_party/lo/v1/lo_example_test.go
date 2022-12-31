package v1

import (
	"fmt"

	"github.com/samber/lo"
)

func ExampleUniq() {
	got := lo.Uniq[string]([]string{"Samuel", "John", "Samuel"})
	fmt.Println(got)

	// Output: [Samuel John]
}

func ExampleFilter() {
	even := lo.Filter[int]([]int{1, 2, 3, 4}, func(x int, index int) bool {
		return x%2 == 0
	})
	fmt.Println(even)

	// Output: [2 4]
}

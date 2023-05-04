package v2s

import (
	"github.com/sourcegraph/conc/iter"
)

// ForEachIdx is the same as ForEach except it also provides the
// index of the element to the callback.
/*func (iter Iterator[T]) ForEachIdx(input []T, f func(int, *T)) {
	......
	var idx atomic.Int64
	// Create the task outside the loop to avoid extra closure allocations.
    // 在 for 循环里创建闭包，传入 idx 参数，然后 wg.Go 去运行。大量闭包，可能会造成 heap 内存增长很快频
    // 繁触发 GC 的性能问题
    // 这里在 for 循环外创建闭包，可以避免产生大量闭包，并通过 atomic 控制 idx，
	task := func() {
		i := int(idx.Add(1) - 1)
		for ; i < numInput; i = int(idx.Add(1) - 1) {
			f(i, &input[i])
		}
	}

	var wg conc.WaitGroup
	for i := 0; i < iter.MaxGoroutines; i++ {
		wg.Go(task)
	}
	wg.Wait()
}*/

func process(values []int) {
	iter.ForEach(values, handle)
}

func handle(e *int) {

}

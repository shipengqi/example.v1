package v2

import "sync"

func mapStream(
	in chan int,
	out chan int,
	f func(int) int,
) {
	tasks := make(chan func())
	taskResults := make(chan chan int)

	// Worker goroutines
	var workerWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for task := range tasks {
				task()
			}
		}()
	}

	// Ordered reader goroutines
	var readerWg sync.WaitGroup
	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		for result := range taskResults {
			item := <-result
			out <- item
		}
	}()

	// Feed the workers with tasks
	for elem := range in {
		resultCh := make(chan int, 1)
		taskResults <- resultCh
		tasks <- func() {
			resultCh <- f(elem)
		}
	}

	// We've exhausted input.
	// Wait for everything to finish
	close(tasks)
	workerWg.Wait()
	close(taskResults)
	readerWg.Wait()
}

package v2

import "sync"

func process(values []int) {
	feeder := make(chan int, 8)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for elem := range feeder {
				handle(elem)
			}
		}()
	}

	for _, value := range values {
		feeder <- value
	}
	close(feeder)
	wg.Wait()
}

func handle(e int) {

}

package main

import "github.com/sourcegraph/conc"

func main() {
	var wg conc.WaitGroup
	wg.Go(doSomethingThatMightPanic)
	// panics with a nice stacktrace
	wg.Wait()
}

func doSomethingThatMightPanic() {
	panic("test panic")
}

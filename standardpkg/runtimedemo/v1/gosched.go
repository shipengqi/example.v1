package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println("----", s)
		runtime.Gosched()
		fmt.Println(s)
	}
}

func main() {
	go say("world")
	say("hello")
}

// Output:
// ---- hello
// ---- world
// hello
// world
// ---- world
// world
// ---- world
// ---- hello
// world
// ---- world
// hello
// ---- hello
// hello
// ---- hello
// world
// ---- world
// hello
// ---- hello
// world
// hello
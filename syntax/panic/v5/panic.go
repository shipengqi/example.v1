package main

import (
	"fmt"
	"runtime"
)

type panicContext struct {
	function string
}

func ProtectRun(entry func()) {
	defer func() {
		err := recover()
		switch err.(type) {
		case runtime.Error:
			fmt.Println(err)
		default:
			fmt.Println("code error:", err)
		}
	}()

	entry()
}

func main() {
	fmt.Println("running")
	ProtectRun(func() {
		fmt.Println("before panic manually")
		panic(&panicContext{"panic manually"})
		fmt.Println("after panic manually")
	})

	ProtectRun(func() {
		fmt.Println("before panic assignment")
		var a *int
		*a = 1
		fmt.Println("after panic assignment")
	})

	fmt.Println("end")
}

// Output:
// running
// before panic manually
// code error: &{panic manually}
// before panic assignment
// runtime error: invalid memory address or nil pointer dereference
// end

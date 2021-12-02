package main

import "fmt"

type Test struct {
	f []string
}

func main() {
	t := Test{}
	fmt.Println(t)
	if t.f == nil {
		fmt.Println("nil")
	}
	if len(t.f) == 0 {
		fmt.Println("zero")
	}
	func1(nil, "test1")
	func1([]string{}, "test2")
	var strs []string
	func1(strs, "test3")
}

func func1(strs []string, name string) {
	if strs == nil {
		fmt.Printf("%s: strs nil\n", name)
	}
	if len(strs) == 0 {
		fmt.Printf("%s: strs zero\n", name)
	}
}

// Output:
// {[]}
// nil
// zero
// test1: strs nil
// test1: strs zero
// test2: strs zero
// test3: strs nil
// test3: strs zero

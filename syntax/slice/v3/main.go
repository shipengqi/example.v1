package main

import "fmt"

func main() {
	var ts []string
	ts2 := make([]string, 0)
	fmt.Println(ts, len(ts))
	fmt.Println(ts2, len(ts2))
	ts2 = append(ts2, "test") // nil slice 只能使用 append 进行操作，否则会 panic
	fmt.Println(ts2[0])
	ts = append(ts, "test2")
	fmt.Println(ts[0])
}

// Output:
// [] 0
// [] 0
// test
// test2

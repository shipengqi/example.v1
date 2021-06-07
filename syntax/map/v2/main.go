package main

import "fmt"

func main() {
	var ts map[string]string
	ts2 := make(map[string]string)
	fmt.Println(ts, len(ts))
	fmt.Println(ts2, len(ts2))
	ts2["test1"] = "test1"
	fmt.Println(ts2, len(ts2))
	// ts["test1"] = "test1" // nil map 不能进行写操作
	v := ts["test1"] // nil map 可以进行读操作
	fmt.Println("-------", v)
}

// Output:
// map[] 0
// map[] 0
// map[test1:test1] 1
// panic: assignment to entry in nil map
//
// goroutine 1 [running]:
// main.main()
//	C:/Code/example.v1/syntax/map/v2/main.go:12 +0x269


// Output:
// map[] 0
// map[] 0
// map[test1:test1] 1
// -------

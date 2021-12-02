package main

import "fmt"

func main() {
	s := make([]int, 0, 10)
	for i := 0; i < 15; i++ {
		s = append(s, i)
	}
	fmt.Println(len(s))
	fmt.Println(cap(s))

	s2 := make([]int, 0, 1024)
	for i := 0; i < 1100; i++ {
		s2 = append(s2, i)
	}
	fmt.Println(len(s2))
	fmt.Println(cap(s2))
}

// Output:
// 15
// 20
// 1100
// 1280 (1024*0.25 + 1024)
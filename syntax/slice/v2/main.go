package main

import "fmt"

func main()  {
	var s []int
	s = append(s, 1)
	fmt.Println(s)

	s2 := new([]int)
	*s2 = append(*s2, 1)
	fmt.Println(s2)
}

// Output
// [1]
// &[1]
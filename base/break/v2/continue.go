package main

import "fmt"

func main() {
Loop:
	for j := 0; j < 3; j++ {
		fmt.Println("loop1", j)
		for a := 0; a < 5; a++ {
			fmt.Println("loop2", a)
			if a > 3 {
				continue Loop
			}
		}
	}
}

// Output:
// loop1 0
// loop2 0
// loop2 1
// loop2 2
// loop2 3
// loop2 4
// loop1 1
// loop2 0
// loop2 1
// loop2 2
// loop2 3
// loop2 4
// loop1 2
// loop2 0
// loop2 1
// loop2 2
// loop2 3
// loop2 4

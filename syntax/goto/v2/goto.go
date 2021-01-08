package main

import "fmt"

func main() {
	fmt.Println("test")
	for a := 0; a < 5; a++ {
		fmt.Println(a)
		if a == 3 {
			goto Loop
		}
	}
Loop:
	fmt.Println("end")
}

// Output：
//test
//0
//1
//2
//3
//end

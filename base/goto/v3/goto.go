package main

import "fmt"

func main() {
	fmt.Println("test")
	goto End
	fmt.Println("test1")
End:
	fmt.Println("end")
}

// Output：
//test
//end

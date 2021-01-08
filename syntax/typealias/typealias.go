package main

import "fmt"

// type alias
type S = string
// type declarations
type S2 string

func main()  {
	var s S = "hello"
	fmt.Println(s)
	var s2 S2 = "world"
	fmt.Println(s2)
	var i interface{}
	i = s2

	switch t := i.(type) {
	case string:
		fmt.Println("string type: ", t)
	case S2:
		fmt.Println("S2 type: ", t)
	default:
		fmt.Println("default")
	}
}

// Output:
//hello
//world
//S2 type:  world

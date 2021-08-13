package main

import "fmt"

type type1 struct{
	a    string
}

type type2 struct{
	a    string
}

func main()  {
	t1 := type1{a: "a"}
	t2 := type1{a: "a"}
	if t1 == t2 {
		fmt.Println("t1 == t2")
	}

	t3 := &type1{a: "a"}
	t4 := &type1{a: "a"}
	if t3 == t4 {
		fmt.Println("t3 == t4")
	}

	if *t3 == *t4 {
		fmt.Println("*t3 == *t4")
	}
}

// Output:
// t1 == t2
// *t3 == *t4
package main

import "fmt"

type TestPoint struct {
	Name string
	Age  int
}

func main() {
	t := &TestPoint{
		Name: "pooky",
		Age:  30,
	}
	merge(t)
	fmt.Printf("%+v\n", t)
	t2 := &TestPoint{
		Name: "pooky",
		Age:  30,
	}
	merge2(t2)
	fmt.Printf("%+v\n", t2)
}

func merge(t *TestPoint) {
	t.Name = "xiaoming"
}

func merge2(t *TestPoint) {
	newT := &TestPoint{
		Name: "xiaoqiang",
		Age:  18,
	}
	t = newT
}


// Output:
// &{Name:xiaoming Age:30}
// &{Name:pooky Age:30}

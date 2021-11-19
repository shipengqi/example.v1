package main

import "fmt"

type TestPoint struct {
	Name string
	Age  int
}

var global *TestPoint
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
	merge3()
	fmt.Printf("%+v\n", global)
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

func merge3() {
	newT := &TestPoint{
		Name: "xiaoqiang",
		Age:  18,
	}
	global = newT
}


// Output:
// &{Name:xiaoming Age:30}
// &{Name:pooky Age:30}

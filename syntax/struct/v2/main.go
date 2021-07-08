package main

import "fmt"

type hello interface {
	Hello()
}

type person struct {
	name string
}

func (p *person) Hello() {
	fmt.Println("Hello, I'm ", p.name)
}

type student struct {
	*person
}

func main() {
	var h hello
	h = &student{
		&person{name: "xiaoming"},
	}
	h.Hello()

	h = &person{name: "xiaoqiang"}
	h.Hello()
}

// Hello, I'm  xiaoming
// Hello, I'm  xiaoqiang

package main

import "fmt"

type hello interface {
	Name() string
	Hello()
}

type person struct {
	name string
}

func (p *person) Hello() {
	fmt.Println("Hello, I'm ", p.Name())
}

func (p *person) Name() string {
	return p.name
}

type student struct {
	*person
}

func (s *student) Hello() {
	fmt.Println("Hello, I'm a student, my name is ", s.Name())
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

// Hello, I'm a student, my name is  xiaoming
// Hello, I'm  xiaoqiang

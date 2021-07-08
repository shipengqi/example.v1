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

func (s *student) Name() string {
	return "student " + s.name
}

func Hello(h hello) {
	fmt.Println("Hello, I'm ", h.Name())
}

func main() {
	var h hello
	h = &student{
		&person{name: "xiaoming"},
	}
	h.Hello() // Hello, I'm  xiaoming, because in go, parent cannot call the child method.
	fmt.Println(h.Name())
	Hello(h)
	h = &person{name: "xiaoqiang"}
	h.Hello()
	fmt.Println(h.Name())
	Hello(h)
}

// Hello, I'm  xiaoming
// student xiaoming
// Hello, I'm  student xiaoming
// Hello, I'm  xiaoqiang
// xiaoqiang
// Hello, I'm  xiaoqiang

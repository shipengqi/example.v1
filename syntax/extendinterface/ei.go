package main

import "fmt"

type Person interface {
	Name() string
	Age() int
}

type Student struct {
	Person

	level int
	score int
	age   int
	name  string
}

func (s *Student) Level() int {
	return s.level
}

func (s *Student) Score() int {
	return s.score
}

//func (s *Student) Name() string {
//	return s.name
//}

//func (s *Student) Age() int {
//	return s.age
//}

func main() {
	var p Person = &Student{
		level: 5,
		score: 100,
		name:  "xiaoming",
	}

	fmt.Println(p.Name())
}

// 可以编译通过
// 但是运行会 runtime panic，因为 Student 并没有实现 Person 接口
// 实现 Person 接口方法集合中的任意方法，都可以执行，例如，只实现 Name 方法，执行 p.Name() 就不会 panic

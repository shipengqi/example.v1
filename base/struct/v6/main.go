package main

import "fmt"

type A struct {
	value1 int
	value2 int
}

func (a A) AddValue(v A) A {
	a.value1 += v.value1
	a.value2 += v.value2
	return a
}

func (a A) AddValue2(v A) {
	a.value1 += v.value1
	a.value2 += v.value2
}

func main() {
	x, z := A{1, 2}, A{1, 2}
	y := A{3, 4}
	x = x.AddValue(y)
	z.value1 += y.value1
	z.value2 += y.value2

	y.AddValue(x)
	fmt.Println(x) // {4 6}
	fmt.Println(z) // {4 6}
	fmt.Println(y) // {3 4}
}

package main

// interface 不是单纯的值，而是分为类型和值
// 必须得类型和值同时都为 nil 的情况下，interface 的 nil 判断才会为 true

import (
	"fmt"
)

type MyErr struct {
	Msg string
}

func main() {
	var e error
	e = GetErr()
	fmt.Println(e == nil) // false
}

func GetErr() *MyErr {
	return nil
}

func (m *MyErr) Error() string {
	return "my error"
}

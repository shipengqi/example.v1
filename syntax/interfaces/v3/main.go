package main

import (
	"fmt"
)

type MyErr struct {
	Msg string
}

func main() {
	var e interface{}
	e = GetErr()
	fmt.Println(e == nil) // true
}

func GetErr() interface{} {
	return nil
}

func (m *MyErr) Error() string {
	return "my error"
}

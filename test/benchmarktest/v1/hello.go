package v1

import (
	"fmt"
	"time"
)

func Hello(name string)  {
	fmt.Println("Hello, ", name)
}

type Big interface {
	Len() int
}

type big struct {
	len int
}

func (b *big) Len() int {
	return b.len
}
func NewBig() Big {
	time.Sleep(time.Second)
	return &big{len: time.Now().Nanosecond()}
}
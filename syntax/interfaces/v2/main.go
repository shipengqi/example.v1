package main

import (
	"fmt"
)

type Base interface {
	do()
}

type App struct {
}

func main() {
	var base Base
	base = GetApp()

	fmt.Println(base)          // <nil>
	fmt.Println(base == nil)   // false
}

func GetApp() *App {
	return nil
}
func (a *App) do() {}

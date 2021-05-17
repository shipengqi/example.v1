package main

import "github.com/shipengqi/example.v1/cli/cmd"

func main() {
	c := cmd.New()
	_ = c.Execute()
}

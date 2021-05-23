package main

import (
	"os"

	"github.com/shipengqi/example.v1/cli/cmd"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

func main() {
	defer recovery()

	c := cmd.New()
	err := c.Execute()
	if err != nil {
		log.Errorf("cmd.Execute(): %v", err)
		os.Exit(ExitCodeError)
	}
	os.Exit(ExitCodeOk)
}

func recovery() {
	if err := recover(); err != nil {
		log.Errorf("[Recovery] panic recovered: %+v", err)
		os.Exit(ExitCodeError)
	}
}
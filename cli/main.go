package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/cmd/certmng"
	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

func main() {
	defer recovery()

	cfg := env.New()
	c := certmng.New(cfg)

	cobra.OnInitialize(func() {
		err := cfg.Init()
		if err != nil {
			panic(errors.Wrap(err, "cfg.Init()"))
		}
		_, err = log.Init(cfg.Log)
		if err != nil {
			panic(errors.Wrap(err, "log.Init()"))
		}
	})

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

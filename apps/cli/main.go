package main

import (
	"github.com/shipengqi/example.v1/apps/cli/cmd/certmng"
	action2 "github.com/shipengqi/example.v1/apps/cli/internal/action"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"github.com/shipengqi/example.v1/apps/cli/pkg/recovery"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

func main() {
	defer recovery.Recovery(ExitCodeError)

	cfg := action2.NewConfiguration()
	c := certmng.New(cfg)

	cobra.OnInitialize(func() {
		err := cfg.Init()
		if err != nil {
			panic(err)
		}
		_, err = log.Init(cfg.Log)
		if err != nil {
			panic(err)
		}
	})

	code := ExitCodeOk
	err := c.Execute()
	if err != nil {
		if err == terminal.InterruptErr {
			log.Warnf("%s.Execute(), interrupted.", c.Name())
		} else if err == action2.DropError {
			log.Warnf("%s.Execute(), exited.", c.Name())
		} else {
			log.Errorf("%s.Execute(): %v", c.Name(), err)
			// If the RunE return error, the PersistentPostRun func will be skipped, so add the following
			// if cfg.Remote == true should skip this
			if !cfg.Remote {
				log.Warn("Additional logging details can be found in:")
				log.Warnf("    %s", log.LogFileName)
			}
			code = ExitCodeError
		}
	}
	os.Exit(code)
}

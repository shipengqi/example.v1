package action

import (
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Check struct {
	name string
	cfg  *Configuration
}

func NewCheck(cfg *Configuration) Interface {
	return &Check{name: "check", cfg: cfg}
}

func (a *Check) Name() string {
	return a.name
}

func (a *Check) PreRun() error {
	log.Info("check certificates.")
	return nil
}

func (a *Check) Run() error {
	return nil
}

func (a *Check) PostRun() error {
	return nil
}

func (a *Check) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}

package action

import (
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Create struct {
	name string
	cfg  *Configuration
}

func NewCreate(cfg *Configuration) Interface {
	return &Create{name: "create", cfg: cfg}
}

func (a *Create) Name() string {
	return a.name
}

func (a *Create) PreRun() error {
	return nil
}

func (a *Create) Run() error {
	log.Info("create certificates.")
	return nil
}

func (a *Create) PostRun() error {
	return nil
}

func (a *Create) Execute() error {
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

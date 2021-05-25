package action

import (
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Renew struct {
	name string
	cfg  *Configuration
}

func NewRenew(cfg *Configuration) Interface {
	return &Renew{name: "renew", cfg: cfg}
}

func (a *Renew) Name() string {
	return a.name
}

func (a *Renew) PreRun() error {
	return nil
}

func (a *Renew) Run() error {
	log.Info("renew certificates.")
	return nil
}

func (a *Renew) PostRun() error {
	return nil
}

func (a *Renew) Execute() error {
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

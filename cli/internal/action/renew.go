package action

import (
	"github.com/pkg/errors"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var DropError = errors.New("Exit")

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
	log.Debug("====================    PRE CHECK    ====================")
	if a.cfg.SkipConfirm {
		return nil
	}
	if a.cfg.Env.RunInPod {
		return nil
	}
	confirm, err := prompt.Confirm("Are you sure to continue")
	if err != nil {
		return err
	}
	if !confirm {
		return DropError
	}
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

package action

import (
	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Renew struct {
	cfg *env.Global
}

func NewRenew(cfg *env.Global) Interface {
	return &Renew{cfg: cfg}
}

func (a *Renew) Name() string {
	return "renew"
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

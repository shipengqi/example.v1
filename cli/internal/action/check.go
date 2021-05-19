package action

import (
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Check struct {
	cfg *config.Global
}

func NewCheck(cfg *config.Global) Interface {
	return &Check{cfg: cfg}
}

func (a *Check) Name() string {
	return "check"
}

func (a *Check) PreRun() error {
	return nil
}

func (a *Check) Run() error {
	log.Info("check certificates.")
	return nil
}

func (a *Check) PostRun() error {
	return nil
}

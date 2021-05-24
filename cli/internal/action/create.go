package action

import (
	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Create struct {
	cfg *env.Global
}

func NewCreate(cfg *env.Global) Interface {
	return &Create{cfg: cfg}
}

func (a *Create) Name() string {
	return "create"
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

package action

import (
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Create struct {
	cfg *config.Global
}

func NewCreate(cfg *config.Global) Interface {
	return &Create{cfg: cfg}
}

func (a *Create) Run() error {
	log.Info("create certificates.")
	return nil
}

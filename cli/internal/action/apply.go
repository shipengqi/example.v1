package action

import (
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Apply struct {
	cfg *config.Global
}

func NewApply(cfg *config.Global) Interface {
	return &Apply{cfg: cfg}
}

func (a *Apply) Run() error {
	log.Info("Start to apply certificates.")
	log.Info("Apply certificates successfully.")
	return nil
}

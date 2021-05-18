package action

import (
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Renew struct {
	cfg *config.Global
}

func NewRenew(cfg *config.Global) Interface {
	return &Renew{cfg: cfg}
}

func (a *Renew) Run() error {
	log.Info("renew certificates.")
	return nil
}

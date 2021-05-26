package action

import (
	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/deployment"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var DropError = errors.New("Exit")

type renew struct {
	*action

	infra  certs.Generator
	deploy certs.Generator
}

func NewRenew(cfg *Configuration) Interface {
	r := &renew{
		action: &action{
			name: "renew",
			cfg:  cfg,
		},
		infra:  infra.New(),
		deploy: deployment.New(),
	}

	return r
}

func (a *renew) Name() string {
	return a.name
}

func (a *renew) PreRun() error {
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

func (a *renew) Run() error {
	log.Debug("====================    RENEW CRT    ====================")
	switch a.cfg.CertType {
	case types.CertTypeInternal:
		break
	case types.CertTypeExternal:
		break
	default:
		return errors.Errorf("unknown cert type: %s", a.cfg.CertType)
	}

	return nil
}

func (a *renew) PostRun() error {
	log.Info("Finished.")
	return nil
}

func (a *renew) Execute() error {
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

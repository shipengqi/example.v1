package action

import (
	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/deployment"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type create struct {
	*action

	infra  certs.Generator
	deploy certs.Generator
}

func NewCreate(cfg *Configuration) Interface {
	c := &create{
		action: &action{
			name: "create",
			cfg:  cfg,
		},
		infra:  infra.New(),
		deploy: deployment.New(),
	}

	return c
}

func (a *create) Name() string {
	return a.name
}

func (a *create) Run() error {
	log.Debug("====================    CREATE CRT    ====================")
	switch a.cfg.NodeType {
	case types.NodeTypeControlPlane:
		break
	case types.NodeTypeWorker:
		break
	default:
		return errors.Errorf("unknown node type: %s", a.cfg.NodeType)
	}
	return nil
}

func (a *create) PostRun() error {
	log.Info("Finished.")
	return nil
}

func (a *create) Execute() error {
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

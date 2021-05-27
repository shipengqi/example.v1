package action

import (
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type apply struct {
	*action
}

func NewApply(cfg *Configuration) Interface {
	return &apply{&action{
		name: applyFlagName,
		cfg:  cfg,
	}}
}

func (a *apply) Name() string {
	return a.name
}

func (a *apply) Run() error {
	log.Debug("====================    APPLY CRT    ====================")
	return sysc.RestartKubeService(a.cfg.Env.CDFNamespace, a.cfg.Env.Version)
}

func (a *apply) PostRun() error {
	log.Info("Apply certificates successfully.")
	return nil
}

func (a *apply) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		log.Warnf("Make sure that you have run the '%s/scripts/renewCert apply' "+
			"on other master nodes.", a.cfg.Env.K8SHome)
		return err
	}
	return a.PostRun()
}

package action

import (
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Apply struct {
	name string
	cfg  *Configuration
}

func NewApply(cfg *Configuration) Interface {
	return &Apply{name: "apply", cfg: cfg}
}

func (a *Apply) Name() string {
	return a.name
}

func (a *Apply) PreRun() error {
	log.Info("start to apply certificates.")
	return nil
}

func (a *Apply) Run() error {
	return sysc.RestartKubeService(a.cfg.Env.CDFNamespace)
}

func (a *Apply) PostRun() error {
	log.Info("Apply certificates successfully.")
	return nil
}

func (a *Apply) Execute() error {
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

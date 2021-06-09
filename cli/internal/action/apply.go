package action

import (
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"strings"
)

const (
	NamespaceKubeSystem = "kube-system"
)

type apply struct {
	*action
}

func NewApply(cfg *Configuration) Interface {
	return &apply{&action{
		name: "apply",
		cfg:  cfg,
	}}
}

func (a *apply) Name() string {
	return a.name
}

func (a *apply) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	return sysc.RestartKubeService(NamespaceKubeSystem, a.cfg.Env.Version)
}

func (a *apply) PostRun() error {
	log.Info("Apply certificates successfully.")
	return nil
}

func (a *apply) Execute() error {
	_ = a.PreRun()
	err := a.Run()
	if err != nil {
		log.Warnf("Make sure that you have run the '%s/scripts/renewCert apply' "+
			"on other master nodes.", a.cfg.Env.K8SHome)
		return err
	}
	return a.PostRun()
}

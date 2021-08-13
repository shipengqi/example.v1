package action

import (
	"github.com/shipengqi/example.v1/apps/cli/internal/sysc"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"strings"
)

const (
	NamespaceKubeSystem = "kube-system"
)

type apply struct {
	*action
}

func NewApply(cfg *Configuration) Interface {
	return &apply{newAction("apply", cfg)}
}

func (a *apply) Name() string {
	return a.name
}

func (a *apply) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	err := sysc.RestartKubeService(NamespaceKubeSystem, a.cfg.Env.Version)
	if err != nil {
		log.Warnf("Make sure that you have run the '%s/scripts/renewCert apply' "+
			"on other master nodes.", a.cfg.Env.K8SHome)
		return err
	}

	log.Info("Apply certificates successfully.")
	return nil
}

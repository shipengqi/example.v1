package action

import (
	"strings"

	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubExternalNotInPod struct {
	*action
}

func NewRenewSubExternalNotInPod(cfg *Configuration) Interface {
	return &renewSubExternalNotInPod{
		action: newAction("renew-sub-external-notinpod", cfg),
	}
}

func (a *renewSubExternalNotInPod) Name() string {
	return a.name
}

func (a *renewSubExternalNotInPod) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	return sysc.RenewRERemoteExecution(a.cfg.Env.CDFNamespace, a.cfg.Namespace,
		a.cfg.Unit, a.cfg.Resource, a.cfg.ResourceField, a.cfg.Cluster.IsPrimary, a.cfg.Validity)
}

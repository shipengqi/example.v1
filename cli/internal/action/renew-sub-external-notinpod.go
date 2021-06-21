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

	err := sysc.RenewRERemoteExecution(a.cfg.Env.CDFNamespace, a.cfg.Namespace,
		a.cfg.Unit, a.cfg.Resource, a.cfg.ResourceField, a.cfg.Cluster.IsPrimary, a.cfg.Validity)
	if err != nil {
		return err
	}

	if !a.cfg.Cluster.IsPrimary {
		return nil
	}
	from := a.cfg.Env.CDFNamespace
	to := a.cfg.Namespace
	if from == to {
		return nil
	}
	log.Debugf("start to copy secrets from: %s, to: %s", from, to)
	secrets := strings.Split(a.cfg.Resource, ",")
	for k := range secrets {
		secret := strings.TrimSpace(secrets[k])
		if len(secret) == 0 {
			continue
		}
		log.Infof("Copying secret: %s ...", secret)
		s, err := a.kube.GetSecret(from, secret)
		if err != nil {
			return err
		}
		_, err = a.kube.ApplySecretBytes(to, secret, s.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

package action

import (
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type check struct {
	*action
}

func NewCheck(cfg *Configuration) Interface {
	return &check{
		action: &action{
			name: "check",
			cfg:  cfg,
		},
	}
}

func (a *check) Name() string {
	return a.name
}

func (a *check) Run() error {
	if len(a.cfg.Cert) > 0 {
		available, err := utils.CheckCrt(a.cfg.Cert)
		if err != nil {
			return err
		}
		if available <= 0 {
			log.Infof("The certificate: %s has already expired.", a.cfg.Cert)
		} else {
			log.Infof("The certificate: %s will expire in %d hour(s).", a.cfg.Cert, available)
		}
	}
	if len(a.cfg.Namespace) > 0 && len(a.cfg.Secret) > 0 {
		client, err := kube.New(a.cfg.Kube)
		if err != nil {
			return err
		}
		secret, err := client.GetSecret(a.cfg.Namespace, a.cfg.Secret)
		if err != nil {
			return err
		}
		for k := range secret.StringData {
			if len(secret.StringData[k]) > 0 {
				available, err := utils.CheckCrtString(secret.StringData[k])
				if err != nil {
					return err
				}
				if available <= 0 {
					log.Infof("The %s.%s has already expired.", a.cfg.Secret, k)
				} else {
					log.Infof("The %s.%s will expire in %d hour(s).", a.cfg.Secret, k, available)
				}
			}
		}
	}

	return nil
}

func (a *check) Execute() error {
	return a.Run()
}

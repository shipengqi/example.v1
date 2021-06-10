package action

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubExternalCustom struct {
	*action

	customCrt []byte
	customKey []byte
}

func NewRenewSubExternalCustom(cfg *Configuration) Interface {
	return &renewSubExternalCustom{
		action: newActionWithKube("renew-sub-external-custom", cfg),
	}
}

func (a *renewSubExternalCustom) Name() string {
	return a.name
}

func (a *renewSubExternalCustom) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	data := make(map[string]string)
	data[a.cfg.ResourceField+".crt"] = string(a.customCrt)
	data[a.cfg.ResourceField+".key"] = string(a.customKey)

	// apply secret
	secrets := strings.Split(a.cfg.Resource, ",")
	for k := range secrets {
		secret := strings.TrimSpace(secrets[k])
		if len(secret) == 0 {
			continue
		}
		log.Debugf("Apply %s in %s", secret, a.cfg.Namespace)
		_, err := a.kube.ApplySecret(a.cfg.Namespace, secret, data)
		if err != nil {
			return errors.Wrapf(err, "apply %s, namespace: %s", secret, a.cfg.Namespace)
		}
	}

	// apply public-ca-certificates configmap
	if len(a.cfg.CACert) > 0 {
		log.Debugf("Read %s", a.cfg.CACert)
		cacertData, err := utils.ReadFile(a.cfg.CACert)
		if err != nil {
			return err
		}
		log.Debugf("Apply %s in %v", ConfigMapNamePublicCA, a.cfg.Namespace)
		newData := make(map[string]string)
		newData["CUS_ca.crt"] = string(cacertData)

		_, err = a.kube.ApplyConfigMap(a.cfg.Namespace, ConfigMapNamePublicCA, newData)
		if err != nil {
			return errors.Wrapf(err, "apply %s, namespace: %s", ConfigMapNamePublicCA, a.cfg.Namespace)
		}
	}

	return nil

}

func (a *renewSubExternalCustom) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	a.cfg.Debug()

	log.Debugf("Read %s", a.cfg.Cert)
	certData, err := utils.ReadFile(a.cfg.Cert)
	if err != nil {
		return err
	}
	crt, err := utils.ParseCrtBytes(certData)
	if err != nil {
		return err
	}
	available := utils.CheckCrtValidity(crt)
	if available <= 0 {
		log.Infof("The certificate: %s has already expired.", a.cfg.Cert)
	} else {
		log.Infof("The certificate: %s will expire in %d hour(s).", a.cfg.Cert, available)
	}

	log.Debugf("Read %s", a.cfg.Key)
	keyData, err := utils.ReadFile(a.cfg.Key)
	if err != nil {
		return err
	}

	a.customCrt = certData
	a.customKey = keyData

	return nil
}

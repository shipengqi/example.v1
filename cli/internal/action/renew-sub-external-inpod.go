package action

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/deployment"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubExternalInPod struct {
	*action

	generator certs.Generator
}

func NewRenewSubExternalInPod(cfg *Configuration) Interface {
	c := &renewSubExternalInPod{
		action: newAction("renew-sub-external-inpod", cfg),
	}

	ns := cfg.Namespace
	if cfg.Cluster.IsPrimary {
		ns = cfg.Env.CDFNamespace
	}
	cfg.Vault.Address = fmt.Sprintf("https://%s.%s:8200", DefaultVaultSvcName, ns)

	cas, err := c.getCAs(cfg.Env.CDFNamespace)
	if err != nil {
		panic(err)
	}
	cfg.Vault.CAs = cas

	g, err := deployment.New(
		cfg.Namespace,
		cfg.Env.CDFNamespace,
		cfg.Cluster.IsPrimary,
		cfg.Kube,
		cfg.Vault)
	if err != nil {
		panic(err)
	}
	c.generator = g

	return c
}

func (a *renewSubExternalInPod) Name() string {
	return a.name
}

func (a *renewSubExternalInPod) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	if !strings.Contains(a.cfg.Resource, SecretNameNginxFrontend) {
		log.Debugf("add secret %s ro resource", SecretNameNginxFrontend)
		a.cfg.Resource = fmt.Sprintf("%s,%s", a.cfg.Resource, SecretNameNginxFrontend)
	}

	return a.generator.GenAndDump(&certs.Certificate{
		CN:       a.cfg.Cluster.ExternalHost,
		UintTime: a.cfg.Unit,
		Validity: a.cfg.Validity,
	}, fmt.Sprintf("%s %s", a.cfg.Resource, a.cfg.ResourceField))

}

func (a *renewSubExternalInPod) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	cm, err := a.kube.GetConfigMap(a.cfg.Env.CDFNamespace, ConfigMapNameCDF)
	if err != nil {
		log.Warnf("kube.GetConfigMap(): %v", err)
	} else {
		a.cfg.Cluster.ExternalHost = cm.Data[ResourceKeyExternalHost]
	}
	a.cfg.Debug()

	return nil
}

func (a *renewSubExternalInPod) getCAs(namespace string) ([][]byte, error) {
	cm, err := a.kube.GetConfigMap(namespace, ConfigMapNamePublicCA)
	if err != nil {
		return nil, err
	}
	ric, ok := cm.Data[ResourceKeyRICCert]
	if !ok {
		return nil, errors.New("RIC ca nil")
	}
	var cas [][]byte
	cas = append(cas, []byte(ric))

	if rid, ok := cm.Data[ResourceKeyRIDCert]; ok {
		log.Debug("got RID ca")
		cas = append(cas, []byte(rid))
	}
	if re, ok := cm.Data[ResourceKeyRECert]; ok {
		log.Debug("got RE ca")
		cas = append(cas, []byte(re))
	}
	if cus, ok := cm.Data[ResourceKeyCUSCert]; ok {
		log.Debug("got CUS ca")
		cas = append(cas, []byte(cus))
	}

	return cas, nil
}

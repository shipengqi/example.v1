package action

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubInternalExpired struct {
	*action

	generator certs.Generator
}

func NewRenewSubInternalExpired(cfg *Configuration) Interface {
	c := &renewSubInternalExpired{
		action: newAction("renew-sub-internal-expired", cfg),
	}

	key, err := c.parseCAKey()
	if err != nil {
		panic(err)
	}

	g, err := infra.New(cfg.CACert, key)
	if err != nil {
		panic(err)
	}
	c.generator = g

	return c
}

func (a *renewSubInternalExpired) Name() string {
	return a.name
}

func (a *renewSubInternalExpired) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	var err error
	log.Info("The certificates have already expired.")
	log.Info("Renew current node certificates")
	err = a.iterate(a.cfg.Host, true, true, a.generator)
	if err != nil {
		return err
	}

	log.Info("Start to apply certificates.")
	err = sysc.RestartKubeService(NamespaceKubeSystem, a.cfg.Env.Version)
	if err != nil {
		log.Warnf("Apply certificates failed.")
		log.Warnf("Make sure that you have run the '%s/scripts/renewCert --renew' on other master nodes.", a.cfg.Env.K8SHome)
		return err
	} else {
		log.Info("Apply certificates on current node successfully.")
	}

	return nil
}

func (a *renewSubInternalExpired) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	if !a.cfg.Env.RunOnMaster {
		return errors.New("You can only renew expired cert on master node")
	}

	hostname, err := sysc.Hostname()
	if err != nil {
		log.Warnf("sysc.Hostname(): %v", err)
	} else {
		a.cfg.Host = hostname
	}
	a.cfg.Host = hostname
	log.Debugf("get local hostname: %s", hostname)

	a.cfg.Debug()

	return nil
}

func (a *renewSubInternalExpired) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	log.Info("Checking the cluster status ...")
	var status error

	for i := 0; i < 60; i++ {
		log.Print(".")
		_, status = a.kube.GetConfigMap(a.cfg.Env.CDFNamespace, ConfigMapNameCDFCluster)
		if status == nil {
			break
		}

		time.Sleep(time.Second)
	}
	log.Print("\n")
	if status != nil {
		log.Warnf("The cluster status is not ready.")
		log.Warnf("Make sure that you have run the '%s/scripts/renewCert --renew' on other master nodes.", a.cfg.Env.K8SHome)
		return errors.Wrap(status, "cluster not ready")
	}

	return nil
}

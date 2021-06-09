package action

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubInternalLocal struct {
	*action

	generator certs.Generator
}

func NewRenewSubInternalLocal(cfg *Configuration) Interface {
	g, err := infra.New(cfg.CACert, cfg.CAKey)
	if err != nil {
		panic(err)
	}

	c := &renewSubInternalLocal{
		action: &action{
			name: "renew-sub-internal-local",
			cfg:  cfg,
		},
		generator: g,
	}

	return c
}

func (a *renewSubInternalLocal) Name() string {
	return a.name
}

func (a *renewSubInternalLocal) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	if a.cfg.Env.RunOnMaster {
		return a.renewMaster()
	}
	return a.renewWorker()
}

func (a *renewSubInternalLocal) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	hostname, err := sysc.Hostname()
	if err != nil {
		return errors.Wrap(err, "get hostname")
	}
	a.cfg.Host = hostname
	log.Debugf("get local hostname: %s", hostname)

	a.cfg.Debug()

	return nil
}

func (a *renewSubInternalLocal) renewWorker() error {
	return nil
}

func (a *renewSubInternalLocal) renewMaster() error {
	return nil
}

func (a *renewSubInternalLocal) renewExistCert() error {

	return nil
}

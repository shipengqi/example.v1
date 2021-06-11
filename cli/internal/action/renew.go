package action

import (
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

var DropError = errors.New("Exit")

type IValidity struct {
	ca     int
	server int
}

type renew struct {
	*action

	iValidity IValidity
}

func NewRenew(cfg *Configuration) Interface {
	return &renew{
		action: newAction("renew", cfg),
	}
}

func (a *renew) Name() string {
	return a.name
}

func (a *renew) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	// ignore checks, if running in a pod
	if a.cfg.Env.RunInPod {
		return nil
	}

	// check CA cert validity
	log.Debugf("Checking %s ...", a.cfg.CACert)
	crt, err := utils.ParseCrt(a.cfg.CACert)
	if err != nil {
		return err
	}
	available := utils.CheckCrtValidity(crt)
	if available <= 0 {
		log.Debugf("The certificate: %s has already expired.", a.cfg.CACert)
		return errors.New("CA certificate expired")
	}
	a.iValidity.ca = available

	// Ignore the following checks, if the --local is true
	if a.cfg.Local {
		return nil
	}

	serverCertPath := path.Join(a.cfg.Env.SSLPath, CertNameKubeletServer+".crt")
	log.Debugf("Checking %s ...", serverCertPath)

	serverCrt, err := utils.ParseCrt(serverCertPath)
	if err != nil {
		return err
	}
	a.iValidity.server = utils.CheckCrtValidity(serverCrt)

	return nil
}

func (a *renew) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	switch a.cfg.CertType {
	case types.CertTypeInternal:
		sub := NewRenewSubInternal(a.cfg, a.iValidity)
		return Execute(sub)
	case types.CertTypeExternal:
		sub := NewRenewSubExternal(a.cfg, a.iValidity)
		return Execute(sub)
	default:
		return errors.Errorf("unknown cert type: %s", a.cfg.CertType)
	}
}

func (a *renew) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	log.Info("Finished!")
	return nil
}

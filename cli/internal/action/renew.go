package action

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var DropError = errors.New("Exit")

type renew struct {
	*action

	expired bool
}

func NewRenew(cfg *Configuration) Interface {
	return &renew{
		action:  newActionWithKube("renew", cfg),
	}
}

func (a *renew) Name() string {
	return a.name
}

func (a *renew) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	if a.cfg.CertType == types.CertTypeExternal {
		return nil
	}
	// create new-certs folder for internal cert
	err := os.MkdirAll(a.cfg.OutputDir, 0744)
	if err != nil {
		return err
	}

	// check cert validity
	crt, err := utils.ParseCrt(a.cfg.Cert)
	if err != nil {
		return err
	}
	available := utils.CheckCrtValidity(crt)
	if available <= 0 {
		log.Infof("The certificate: %s has already expired.", a.cfg.Cert)
	} else {
		log.Infof("The certificate: %s will expire in %d hour(s).", a.cfg.Cert, available)
	}

	// confirm
	if a.cfg.SkipConfirm {
		return nil
	}
	if a.cfg.Env.RunInPod {
		return nil
	}
	confirm, err := prompt.Confirm("Are you sure to continue")
	if err != nil {
		return err
	}
	if !confirm {
		return DropError
	}
	return nil
}

func (a *renew) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	log.Info("Renewing certificates ...")

	switch a.cfg.CertType {
	case types.CertTypeInternal:
		log.Infof("Cert type: %s", types.CertTypeInternal)
		return a.renewInternal()
	case types.CertTypeExternal:
		log.Infof("Cert type: %s", types.CertTypeExternal)
		return a.renewExternal()
	default:
		return errors.Errorf("unknown cert type: %s", a.cfg.CertType)
	}
}

func (a *renew) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	log.Info("Finished!")
	return nil
}

func (a *renew) renewExternal() error {
	var sub Interface
	if len(a.cfg.Cert) > 0 && len(a.cfg.Key) > 0 {
		sub = NewRenewSubExternalCustom(a.cfg)
	}

	if !a.cfg.Env.RunInPod {
		sub = NewRenewSubExternalNotInPod(a.cfg)
	} else {
		sub = NewRenewSubExternalInPod(a.cfg)
	}
	return sub.Execute()
}

func (a *renew) renewInternal() error {
	if a.cfg.Local {
		sub := NewRenewSubInternalLocal(a.cfg)
		return sub.Execute()
	}
	if a.expired {
		sub := NewRenewSubInternalExpired(a.cfg)
		err := sub.Execute()
		if err != nil {
			return err
		}
	}
	sub := NewRenewSubInternalAvailable(a.cfg)
	return sub.Execute()
}

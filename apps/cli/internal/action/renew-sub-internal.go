package action

import (
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"github.com/shipengqi/example.v1/apps/cli/pkg/prompt"
	"os"
	"strings"
)

type renewSubInternal struct {
	*action

	iValidity IValidity
}

func NewRenewSubInternal(cfg *Configuration, validity IValidity) Interface {
	return &renewSubInternal{
		action:    newAction("renew-sub-internal", cfg),
		iValidity: validity,
	}
}

func (a *renewSubInternal) Name() string {
	return a.name
}

func (a *renewSubInternal) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	if !a.cfg.SkipConfirm {
		confirm, err := prompt.Confirm("Are you sure to continue")
		if err != nil {
			return err
		}
		if !confirm {
			return DropError
		}
	}

	if a.cfg.Local {
		log.Info("Renewing internal certificates locally ...")
		sub := NewRenewSubInternalLocal(a.cfg)
		return Execute(sub)
	}
	if a.iValidity.server <= 0 {
		log.Info("Renewing expired internal certificates ...")
		sub := NewRenewSubInternalExpired(a.cfg)
		err := Execute(sub)
		if err != nil {
			return err
		}
	}
	log.Info("Renewing available internal certificates ...")
	sub := NewRenewSubInternalAvailable(a.cfg)
	return Execute(sub)
}

func (a *renewSubInternal) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	days := a.iValidity.ca / 24
	if days < a.cfg.Validity {
		log.Warnf("The internal root CA certificate on the current node "+
			"will expire in %d day(s).", days)
		log.Warnf("The certificate validity period must less than %d.", days)
	}

	// Ignore the following checks, if the --local is true
	if !a.cfg.Local {
		if a.iValidity.server <= 0 {
			log.Info("The internal certificates have already expired.")
		} else {
			log.Infof("The internal certificates will expire in %d hour(s).", a.iValidity.server)
		}
	}

	// create new-certs folder for internal cert
	err := os.MkdirAll(a.cfg.OutputDir, 0744)
	if err != nil {
		return err
	}

	return nil
}

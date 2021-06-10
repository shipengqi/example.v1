package action

import (
	"os"
	"path"
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
	} else {
		days := available/24
		log.Warnf("The internal root CA certificate on the current node "+
			"will expire in %d day(s).", days)
		log.Warnf("The certificate validity period must less than %d.", days)
	}

	serverCertPath := path.Join(a.cfg.Env.SSLPath, CertNameKubeletServer+".crt")
	log.Debugf("Checking %s ...", serverCertPath)

	serverCrt, err := utils.ParseCrt(serverCertPath)
	if err != nil {
		return err
	}
	available = utils.CheckCrtValidity(serverCrt)
	if available <= 0 {
		log.Info("The internal certificates have already expired.")
	} else {
		log.Infof("The internal certificates will expire in %d hour(s).", available)
		a.expired = true
	}

	// create new-certs folder for internal cert
	if a.cfg.CertType == types.CertTypeInternal {
		err := os.MkdirAll(a.cfg.OutputDir, 0744)
		if err != nil {
			return err
		}
	}
	if a.cfg.CertType == types.CertTypeExternal {
		if a.expired {
			log.Error("The internal certificates have already expired.")
			log.Errorf("You should run the '%s/scripts/renewCert --renew' to "+
				"renew the internal certificates firstly.", a.cfg.Env.K8SHome)
			return errors.New("internal certificates expired")
		}

		log.Debug("Checking external CA certificate ...")
	}
	// ignore ask in pod
	if a.cfg.Env.RunInPod {
		return nil
	}

	// confirm
	if a.cfg.SkipConfirm {
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
	return Execute(sub)
}

func (a *renew) renewInternal() error {
	if a.cfg.Local {
		sub := NewRenewSubInternalLocal(a.cfg)
		return Execute(sub)
	}
	if a.expired {
		sub := NewRenewSubInternalExpired(a.cfg)
		err := Execute(sub)
		if err != nil {
			return err
		}
	}
	sub := NewRenewSubInternalAvailable(a.cfg)
	return Execute(sub)
}

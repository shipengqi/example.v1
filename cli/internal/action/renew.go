package action

import (
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"os"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/deployment"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var DropError = errors.New("Exit")

type renew struct {
	*action

	generator certs.Generator
	expired   bool
}

func NewRenew(cfg *Configuration) Interface {
	var (
		g   certs.Generator
		err error
	)

	if cfg.CertType == types.CertTypeExternal {
		g, err = deployment.New(cfg.Namespace, cfg.Kube, cfg.Vault)
	} else {
		g, err = infra.New(cfg.CACert, cfg.CAKey)
	}
	if err != nil {
		panic(err)
	}

	return &renew{
		action: &action{
			name: "renew",
			cfg:  cfg,
		},
		generator: g,
	}
}

func (a *renew) Name() string {
	return a.name
}

func (a *renew) PreRun() error {
	log.Debug("*****  RENEW PRE RUN  *****")
	a.cfg.Debug()

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
	log.Debug("*****  RENEW CRT  *****")
	switch a.cfg.CertType {
	case types.CertTypeInternal:
		log.Debugf("cert type: %s", types.CertTypeInternal)
		return a.generator.GenAndDump(&certs.Certificate{
			CN:       a.cfg.Host,
			UintTime: a.cfg.Unit,
			Validity: a.cfg.Validity,
		}, a.cfg.OutputDir)
	case types.CertTypeExternal:
		log.Debugf("cert type: %s", types.CertTypeExternal)
		return a.generator.GenAndDump(&certs.Certificate{
			CN:       a.cfg.Host,
			UintTime: a.cfg.Unit,
			Validity: a.cfg.Validity,
		}, a.cfg.OutputDir)
	default:
		return errors.Errorf("unknown cert type: %s", a.cfg.CertType)
	}
}

func (a *renew) PostRun() error {
	log.Info("Finished.")
	return nil
}

func (a *renew) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}

func (a *renew) renewExternal() error {
	return nil
}

func (a *renew) renewInternal() error {
	return nil
}

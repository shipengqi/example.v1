package action

import (
	"github.com/shipengqi/example.v1/apps/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/apps/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/apps/cli/internal/types"
	"github.com/shipengqi/example.v1/apps/cli/internal/utils"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type create struct {
	*action

	generator certs.Generator
}

func NewCreate(cfg *Configuration) Interface {
	c := &create{
		action: newActionWithoutKube("create", cfg),
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

func (a *create) Name() string {
	return a.name
}

func (a *create) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	log.Info("Creating certificates ...")
	var master bool

	switch a.cfg.NodeType {
	case types.NodeTypeControlPlane:
		master = true
		break
	case types.NodeTypeWorker:
		master = false
		break
	default:
		return errors.Errorf("unknown node type: %s", a.cfg.NodeType)
	}
	return a.iterate(a.cfg.Host, master, false, a.generator)
}

func (a *create) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	a.cfg.Debug()

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

	days := available / 24
	if days < a.cfg.Validity {
		log.Warnf("The internal root CA certificate on the current node "+
			"will expire in %d day(s).", days)
		log.Warnf("The certificate validity period must less than %d.", days)
	}

	// create new-certs folder for internal cert
	return os.MkdirAll(a.cfg.OutputDir, 0744)
}

func (a *create) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	log.Info("Finished!")
	return nil
}

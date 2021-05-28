package action

import (
	"net"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type create struct {
	*action

	generator certs.Generator
}

func NewCreate(cfg *Configuration) Interface {
	g, err := infra.New(cfg.CACert, cfg.CAKey)
	if err != nil {
		panic(err)
	}
	c := &create{
		action: &action{
			name: "create",
			cfg:  cfg,
		},
		generator: g,
	}

	return c
}

func (a *create) Name() string {
	return a.name
}

func (a *create) Run() error {
	log.Debug("*****  CREATE CRT  *****")
	var isMater bool

	switch a.cfg.NodeType {
	case types.NodeTypeControlPlane:
		isMater = true
		break
	case types.NodeTypeWorker:
		isMater = false
		break
	default:
		return errors.Errorf("unknown node type: %s", a.cfg.NodeType)
	}

	for _, v := range CertificateSet {
		if !v.CanDep(isMater) {
			continue
		}

		dns := make([]string, 0)
		ips := make([]net.IP, 0)
		if v.IsServerCert() {
			var sanSvcIp string
			log.Debugf("server cert: %s", v.Name)
			d, i, s := parseSan(a.cfg.ServerCertSan)
			if d != nil {
				dns = append(dns, d ...)
			}
			if i != nil {
				ips = append(ips, i ...)
			}
			if len(s) > 0 {
				sanSvcIp = s
			}
			v.SetIPs(ips, sanSvcIp)
			v.SetDNS(dns)
		}
		v.SetCN("address")
		return a.generator.GenAndDump(v.Certificate, "")
	}

	return nil
}

func (a *create) PostRun() error {
	log.Info("Finished.")
	return nil
}

func (a *create) Execute() error {
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

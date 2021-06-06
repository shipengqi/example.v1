package action

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type check struct {
	*action
}

func NewCheck(cfg *Configuration) Interface {
	return &check{
		action: &action{
			name: "check",
			cfg:  cfg,
		},
	}
}

func (a *check) Name() string {
	return a.name
}

func (a *check) PreRun() error {
	log.Debug("*****  CHECK PRE RUN  *****")
	a.cfg.Debug()

	if len(a.cfg.Cert) == 0 && len(a.cfg.Resource) == 0 {
		return errors.Errorf("Please ")
	}

	if len(a.cfg.Resource) > 0 {
		if len(a.cfg.ResourceField) == 0 {
			return errors.Errorf("Flag --field is reqiured!")
		}
		if len(a.cfg.Namespace) == 0 {
			return errors.Errorf("Flag --namespace is reqiured!")
		}
	}

	return nil
}

func (a *check) Run() error {
	if len(a.cfg.Cert) > 0 {
		log.Debugf("check cert: %s", a.cfg.Cert)
		crt, err := utils.ParseCrt(a.cfg.Cert)
		if err != nil {
			return err
		}
		printCrtInfo(crt)
	}
	if len(a.cfg.Resource) > 0 {
		var certStr string
		var ok bool
		sType, sName := parseSourceName(a.cfg.Resource)
		log.Debugf("check cert %s: %s, namespace: %s", sType, sName, a.cfg.Namespace)
		client, err := kube.New(a.cfg.Kube)
		if err != nil {
			return err
		}
		if sType == types.SourceTypeConfigMap {
			cm, err := client.GetConfigMap(a.cfg.Namespace, sName)
			if err != nil {
				return err
			}
			log.Debugf("check cert %s.%s", sName, a.cfg.ResourceField)
			certStr, ok = cm.Data[a.cfg.ResourceField]
			if !ok {
				return errors.Errorf("field: %s not found", a.cfg.ResourceField)
			}
		} else if sType == types.SourceTypeSecret {
			secret, err := client.GetSecret(a.cfg.Namespace, sName)
			if err != nil {
				return err
			}
			log.Debugf("check cert %s.%s", sName, a.cfg.ResourceField)
			var certBytes []byte
			certBytes, ok = secret.Data[a.cfg.ResourceField]
			if !ok {
				return errors.Errorf("field: %s not found", a.cfg.ResourceField)
			}
			certStr = string(certBytes)
		} else {
			return errors.Errorf("unknown source type: %s", sType)
		}

		crt, err := utils.ParseCrtString(certStr)
		if err != nil {
			return err
		}

		printCrtInfo(crt)
	}

	return nil
}

func (a *check) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	return a.Run()
}

func parseSourceName(source string) (sType, sName string) {
	names := strings.SplitN(source, ".", 2)
	if len(names) >= 2 {
		return names[0], names[1]
	}
	if len(names) == 1 {
		return names[0], ""
	}

	return "", ""
}

func printCrtInfo(cert *x509.Certificate) {
	log.Info("Certificate Information:")
	log.Infof("  Issuer: %s", cert.Issuer)
	log.Infof("  NotBefore: %s", cert.NotBefore.String())
	log.Infof("  NotAfter: %s", cert.NotAfter.String())
	log.Infof("  Subject: %s", cert.Subject)

	dnsStr := strings.Join(cert.DNSNames, ",")
	ipBuf := new(bytes.Buffer)

	for k := range cert.IPAddresses {
		if k == 0 {
			_, _ = fmt.Fprintf(ipBuf, "%s", cert.IPAddresses[k].String())
		} else {
			_, _ = fmt.Fprintf(ipBuf, ", %s", cert.IPAddresses[k].String())
		}
	}
	log.Infof("  DNSNames: %s", dnsStr)
	log.Infof("  IPAddresses: %s", ipBuf.String())
	log.Infof("  KeyUsage: %v", cert.KeyUsage)
	log.Infof("  ExtKeyUsage: %v", cert.ExtKeyUsage)
	log.Infof("  IsCA: %v", cert.IsCA)
	log.Info("")
	log.Info("Certificate Validity:")
	available := utils.CheckCrtValidity(cert)
	if available <= 0 {
		log.Infof("  The certificate has already expired.")
	} else {
		log.Infof("  The certificate will expire in %d hour(s).", available)
	}
	log.Info("")
}

package action

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"

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

func (a *check) Run() error {
	if len(a.cfg.Cert) > 0 {
		log.Debugf("check cert: %s", a.cfg.Cert)
		crt, err := utils.ParseCrt(a.cfg.Cert)
		if err != nil {
			return err
		}
		printCrtInfo(crt)
	}
	if len(a.cfg.Namespace) > 0 && len(a.cfg.Secret) > 0 {
		log.Debugf("check cert secret: %s, namespace: %s", a.cfg.Secret, a.cfg.Namespace)
		client, err := kube.New(a.cfg.Kube)
		if err != nil {
			return err
		}
		secret, err := client.GetSecret(a.cfg.Namespace, a.cfg.Secret)
		if err != nil {
			return err
		}
		for k := range secret.StringData {
			if len(secret.StringData[k]) > 0 {
				crt, err := utils.ParseCrtString(secret.StringData[k])
				if err != nil {
					return err
				}

				available := utils.CheckCrtStringValidity(crt)
				if available <= 0 {
					log.Infof("The %s.%s has already expired.", a.cfg.Secret, k)
				} else {
					log.Infof("The %s.%s will expire in %d hour(s).", a.cfg.Secret, k, available)
				}
				log.Info("")
				printCrtInfo(crt)
			}
		}
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

func printCrtInfo(cert *x509.Certificate) {

	log.Infof("Issuer: %s", cert.Issuer)
	log.Infof("NotBefore: %s", cert.NotBefore.String())
	log.Infof("NotAfter: %s", cert.NotAfter.String())
	log.Infof("Subject: %s", cert.Subject)

	dnsStr := strings.Join(cert.DNSNames, ",")
	ipBuf := new(bytes.Buffer)

	for k := range cert.IPAddresses {
		if k == 0 {
			_, _ = fmt.Fprintf(ipBuf, "%s",  cert.IPAddresses[k].String())
		} else {
			_, _ = fmt.Fprintf(ipBuf, ", %s",  cert.IPAddresses[k].String())
		}
	}
	log.Infof("DNSNames: %s", dnsStr)
	log.Infof("IPAddresses: %s", ipBuf.String())
	log.Infof("KeyUsage: %v", cert.KeyUsage)
	log.Infof("ExtKeyUsage: %v", cert.ExtKeyUsage)
	log.Infof("IsCA: %v", cert.IsCA)

	available := utils.CheckCrtValidity(cert)
	if available <= 0 {
		log.Infof("The certificate has already expired.")
	} else {
		log.Infof("The certificate will expire in %d hour(s).", available)
	}

}

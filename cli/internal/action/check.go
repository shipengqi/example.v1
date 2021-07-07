package action

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type check struct {
	*action
}

func NewCheck(cfg *Configuration) Interface {
	return &check{
		action: newAction("check", cfg),
	}
}

func (a *check) Name() string {
	return a.name
}

func (a *check) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))
	a.cfg.Debug()

	if len(a.cfg.Cert) == 0 && len(a.cfg.Resource) == 0 {
		return errors.Errorf("Please specify a resource to check, use flag '--tls-cert' or '--resource'.")
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
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

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

		if sType == types.SourceTypeConfigMap {
			cm, err := a.kube.GetConfigMap(a.cfg.Namespace, sName)
			if err != nil {
				return err
			}
			log.Debugf("check cert %s.%s", sName, a.cfg.ResourceField)
			certStr, ok = cm.Data[a.cfg.ResourceField]
			if !ok {
				return errors.Errorf("field: %s not found", a.cfg.ResourceField)
			}
		} else if sType == types.SourceTypeSecret {
			secret, err := a.kube.GetSecret(a.cfg.Namespace, sName)
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

func parseSourceName(source string) (sType, sName string) {
	names := strings.SplitN(source, "/", 2)
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
	log.Infof("  KeyUsage: %s", keyUsageToString(cert.KeyUsage))
	log.Infof("  ExtKeyUsage: %s", extKeyUsageToString(cert.ExtKeyUsage))
	log.Infof("  SignatureAlgorithm: %s", cert.SignatureAlgorithm.String())
	log.Infof("  PublicKeyAlgorithm: %s", cert.PublicKeyAlgorithm.String())
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

type keyUsageDetail struct {
	keyUsage x509.KeyUsage
	name     string
}

type extKeyUsageDetail struct {
	extKeyUsage x509.ExtKeyUsage
	name        string
}

var keyUsageDetails = [...]keyUsageDetail{
	{x509.KeyUsageDigitalSignature, "DigitalSignature"},
	{x509.KeyUsageContentCommitment, "ContentCommitment"},
	{x509.KeyUsageKeyEncipherment, "KeyEncipherment"},
	{x509.KeyUsageDataEncipherment, "DataEncipherment"},
	{x509.KeyUsageKeyAgreement, "KeyAgreement"},
	{x509.KeyUsageCertSign, "CertSign"},
	{x509.KeyUsageCRLSign, "CRLSign"},
	{x509.KeyUsageEncipherOnly, "EncipherOnly"},
	{x509.KeyUsageDecipherOnly, "DecipherOnly"},
}

var extKeyUsageDetails = [...]extKeyUsageDetail{
	{x509.ExtKeyUsageAny, "Any"},
	{x509.ExtKeyUsageServerAuth, "ServerAuth"},
	{x509.ExtKeyUsageClientAuth, "ClientAuth"},
	{x509.ExtKeyUsageCodeSigning, "CodeSigning"},
	{x509.ExtKeyUsageEmailProtection, "EmailProtection"},
	{x509.ExtKeyUsageIPSECEndSystem, "IPSECEndSystem"},
	{x509.ExtKeyUsageIPSECTunnel, "IPSECTunnel"},
	{x509.ExtKeyUsageIPSECUser, "IPSECUser"},
	{x509.ExtKeyUsageTimeStamping, "TimeStamping"},
	{x509.ExtKeyUsageOCSPSigning, "OCSPSigning"},
	{x509.ExtKeyUsageMicrosoftServerGatedCrypto, "MicrosoftServerGatedCrypto"},
	{x509.ExtKeyUsageNetscapeServerGatedCrypto, "NetscapeServerGatedCrypto"},
	{x509.ExtKeyUsageMicrosoftCommercialCodeSigning, "MicrosoftCommercialCodeSigning"},
	{x509.ExtKeyUsageMicrosoftKernelCodeSigning, "MicrosoftKernelCodeSigning"},
}

func keyUsageToString(keyUsage x509.KeyUsage) string {
	for _, details := range keyUsageDetails {
		if details.keyUsage == keyUsage {
			return details.name
		}
	}
	return ""
}

func extKeyUsageToString(extKeyUsages []x509.ExtKeyUsage) string {
	var tmp []string
	for _, details := range extKeyUsageDetails {
		for k := range extKeyUsages {
			if details.extKeyUsage == extKeyUsages[k] {
				tmp = append(tmp, details.name)
			}
		}
	}
	if len(tmp) == 0 {
		return ""
	}
	return strings.Join(tmp, ",")
}

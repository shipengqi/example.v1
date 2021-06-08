package action

import (
	"crypto/x509"
	"fmt"
	"net"
	"strings"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ServerCertSuffix      = "server"
	KubeSpiServerCertName = "kube-api-server"
	KubeletClientCertName = "kubelet-kube-api-client"
	KubeletCertCNPrefix   = "system:node:"
)


type CertificateSetItem struct {
	*certs.Certificate

	Secret string
	Deploy int
}

func (i *CertificateSetItem) CanDep(isMaster bool) bool {
	if i.Deploy == types.DepMasterAndWorker {
		return true
	}
	if isMaster && i.Deploy == types.DepMaster {
		return true
	}
	if !isMaster && i.Deploy == types.DepWorker {
		return true
	}
	return false
}

func (i *CertificateSetItem) CombineServerSan(dns []string, ips []net.IP, cn, san, defaultSvcIp string) {
	if i.IsKubeletClientCert() {
		cn = fmt.Sprintf("%s%s", KubeletCertCNPrefix, cn)
	}
	if i.IsServerCert() {
		log.Debugf("server cert: %s", i.Name)
		serverDNSNames, serverIps, svcIp := parseSan(san)
		if serverDNSNames != nil {
			dns = append(dns, serverDNSNames ...)
		}
		if serverIps != nil {
			ips = append(ips, serverIps ...)
		}
		if len(svcIp) > 0 {
			ips = append(ips, svcIp)
		} else if len(defaultSvcIp) > 0 {
			defaultIp := net.ParseIP(defaultSvcIp)
			ips = append(ips, defaultIp)
		}

		if i.IsK8sApiServerCert() {
			dns = append(
				dns,
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				"kubernetes.default.svc.cluster.local",
			)
		}
	}

	i.CN = cn
	i.DNSNames = dns
	i.IPs = ips
}

func (i *CertificateSetItem) IsServerCert() bool {
	return strings.HasSuffix(i.Name, ServerCertSuffix)
}

func (i *CertificateSetItem) IsK8sApiServerCert() bool {
	return i.Name == KubeSpiServerCertName
}

func (i *CertificateSetItem) IsKubeletClientCert() bool {
	return i.Name == KubeletClientCertName
}

var CertificateSet = []CertificateSetItem{
	{
		Certificate: &certs.Certificate{
			Name:   "etcd-server",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "common-etcd-client",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-api-etcd-client",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-api-kubelet-client",
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:kubelet-api-admin"},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-api-proxy-client",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-api-server",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-controller-kube-api-client",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-controller-manager",
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kubectl-kube-api-client",
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:masters"},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kubelet-kube-api-client",
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:nodes"},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kubelet-server",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-scheduler-kube-api-client",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-scheduler",
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "metrics-server",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			CN:           "metrics-server.kube-system",
		},
		Secret: "metrics-server-cert.kube-system",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:   "kube-registry",
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Secret: "kube-registry-cert.<namespace>",
		Deploy: types.DepMaster,
	},
}

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

	Name   string
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

func (i *CertificateSetItem) CombineServerSan(dns []string, ips []net.IP, cn, san string) {
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
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "etcd-server",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "common-etcd-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-etcd-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:kubelet-api-admin"},
		},
		Name:   "kube-api-kubelet-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-proxy-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-api-server",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-controller-manager",
		},
		Name:   "kube-controller-kube-api-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:masters"},
		},
		Name:   "kubectl-kube-api-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:nodes"},
		},
		Name:   "kubelet-kube-api-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kubelet-server",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-scheduler",
		},
		Name:   "kube-scheduler-kube-api-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			CN:           "metrics-server.kube-system",
		},
		Name:   "metrics-server",
		Secret: "metrics-server-cert.kube-system",
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-registry",
		Secret: "kube-registry-cert.<namespace>",
		Deploy: types.DepMaster,
	},
}

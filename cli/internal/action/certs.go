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
	ServerCertSuffix         = "server"

	KubeletCertCNPrefix      = "system:node:"
	KubeSchedulerCN          = "system:kube-scheduler"
	KubeControllerClientCN   = "system:kube-controller-manager"

	OrgKubeletAdmin          = "system:kubelet-api-admin"
	OrgKubectlMaster         = "system:masters"
	OrgKubectlNode           = "system:nodes"

	CertNameKubeApiServer    = "kube-api-server"
	CertNameKubeletClient    = "kubelet-kube-api-client"
	CertNameKubeletAPIClient = "kube-api-kubelet-client"
	CertNameKubeletServer    = "kubelet-server"
	CertNameKubeRegistry     = "kube-registry"
	CertNameMetricServer     = "metrics-server"
	CertNameSchedulerClient  = "kube-scheduler-kube-api-client"
	CertNameETCDClient       = "common-etcd-client"
	CertNameETCDAPIClient    = "kube-api-etcd-client"
	CertNameETCDServer       = "etcd-server"
	CertNameKubectlClient    = "kubectl-kube-api-client"
	CertNameControllerClient = "kube-controller-kube-api-client"
	CertNameProxyClient      = "kube-api-proxy-client"
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
			dns = append(dns, serverDNSNames...)
		}
		if serverIps != nil {
			ips = append(ips, serverIps...)
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
	return i.Name == CertNameKubeApiServer
}

func (i *CertificateSetItem) IsKubeletClientCert() bool {
	return i.Name == CertNameKubeletClient
}

func (i *CertificateSetItem) IsKubeRegistryCert() bool {
	return i.Name == CertNameKubeRegistry
}

var CertificateSet = []CertificateSetItem{
	{
		Certificate: &certs.Certificate{
			Name:         CertNameETCDServer,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameETCDClient,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameETCDAPIClient,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:          CertNameKubeletAPIClient,
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{OrgKubeletAdmin},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameProxyClient,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameKubeApiServer,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameControllerClient,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           KubeControllerClientCN,
		},
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:          CertNameKubectlClient,
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{OrgKubectlMaster},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:          CertNameKubeletClient,
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{OrgKubectlNode},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameKubeletServer,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameSchedulerClient,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           KubeSchedulerCN,
		},
		Deploy: types.DepMaster,
	},
}

var CertificateSecretSet = []CertificateSetItem{
	{
		Certificate: &certs.Certificate{
			Name:         CertNameMetricServer,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			CN:           fmt.Sprintf("%s.%s", CertNameMetricServer, NamespaceKubeSystem),
		},
		Secret: fmt.Sprintf("%s-cert.%s", CertNameMetricServer, NamespaceKubeSystem),
		Deploy: types.DepMaster,
	},
	{
		Certificate: &certs.Certificate{
			Name:         CertNameKubeRegistry,
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Secret: fmt.Sprintf("%s-cert", CertNameKubeRegistry),
		Deploy: types.DepMaster,
	},
}

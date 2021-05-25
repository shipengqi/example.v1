package action

import (
	"crypto/x509"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
)

const (
	DepWorker int = iota
	DepMaster
	DepMasterAndWorker
)

type CertificateSetItem struct {
	certs.Certificate

	Name   string
	Secret string
	Deploy int
}

var CertificateSet = []CertificateSetItem{
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "etcd-server",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "common-etcd-client",
		Deploy: DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-etcd-client",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:kubelet-api-admin"},
		},
		Name:   "kube-api-kubelet-client",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-proxy-client",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-api-server",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-controller-manager",
		},
		Name:   "kube-controller-kube-api-client",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:masters"},
		},
		Name:   "kubectl-kube-api-client",
		Deploy: DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:nodes"},
		},
		Name:   "kubelet-kube-api-client",
		Deploy: DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kubelet-server",
		Deploy: DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-scheduler",
		},
		Name:   "kube-scheduler-kube-api-client",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			CN:           "metrics-server.kube-system",
		},
		Name:   "metrics-server",
		Secret: "metrics-server-cert.kube-system",
		Deploy: DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-registry",
		Secret: "kube-registry-cert.<namespace>",
		Deploy: DepMaster,
	},
}

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Execute() error
}

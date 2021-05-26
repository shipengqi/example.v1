package action

import (
	"crypto/x509"
	"strings"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type CertificateSetItem struct {
	certs.Certificate

	Name   string
	Secret string
	Deploy int
}

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Execute() error
}

var CertificateSet = []CertificateSetItem{
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "etcd-server",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "common-etcd-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-etcd-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:kubelet-api-admin"},
		},
		Name:   "kube-api-kubelet-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		Name:   "kube-api-proxy-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-api-server",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-controller-manager",
		},
		Name:   "kube-controller-kube-api-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:masters"},
		},
		Name:   "kubectl-kube-api-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:      x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Organizations: []string{"system:nodes"},
		},
		Name:   "kubelet-kube-api-client",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		Name:   "kubelet-server",
		Deploy: types.DepMasterAndWorker,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			CN:           "system:kube-scheduler",
		},
		Name:   "kube-scheduler-kube-api-client",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			CN:           "metrics-server.kube-system",
		},
		Name:   "metrics-server",
		Secret: "metrics-server-cert.kube-system",
		Deploy: types.DepMaster,
	},
	{
		Certificate: certs.Certificate{
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			ExtKeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		},
		Name:   "kube-registry",
		Secret: "kube-registry-cert.<namespace>",
		Deploy: types.DepMaster,
	},
}

type action struct {
	name string
	cfg *Configuration
}

func (a *action) Name() string {
	return "[action]"
}

func (a *action) PreRun() error {
	log.Debugf("====================    %s PreRun    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) Run() error {
	log.Debugf("====================   %s Run    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) PostRun() error {
	log.Debugf("====================   %s PostRun    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) Execute() error {
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

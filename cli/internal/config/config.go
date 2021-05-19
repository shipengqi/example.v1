package config

import (
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type Global struct {
	Version            string
	K8SHome            string
	CDFNamespace       string
	RuntimeCDFDataHome string
	SSLPath            string
	NewCertPath        string
	CertRole           string
	VaultAddress       string
	CertType           string
	Username           string
	Password           string
	SSHKey             string
	Cert               string
	Key                string
	CACert             string
	Namespace          string
	LogLevel           string
	LogOutput          string
	Unit               string
	KubeConfig         string
	CAKey              string
	NodeType           string
	Host               string
	OutputDir          string
	ServerCertSan      string
	Install            bool
	Apply              bool
	Renew              bool
	SkipConfirm        bool
	Remote             bool
	Local              bool
	InContainer        bool
	Period             int
	Kube               *kube.Config
	Log                *log.Config
	Vault              *vault.Config
}

func (g *Global) Init() error {
	return nil
}

func combine() *Global {
	return nil
}

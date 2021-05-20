package config

import (
	"github.com/shipengqi/example.v1/cli/internal/flags"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type Global struct {
	Flags              *flags.Global
	Env                *Envs
	Kube               *kube.Config
	Log                *log.Config
	Vault              *vault.Config

	Version            string
	SSLPath            string
	CertRole           string
	VaultAddress       string
	LogLevel           string
	LogOutput          string
}

func (g *Global) Init() error {
	return nil
}

func combine() *Global {
	return nil
}



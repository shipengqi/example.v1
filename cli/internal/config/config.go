package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/flags"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

const (
	DefaultVaultAddr = "https://127.0.0.1:8200"
	DefaultVaultRole = "coretech"
	EnvKeyVaultRole  = "CERTIFICATE_ROLE"
)

type Global struct {
	*flags.Global

	Env   *Envs
	Log   *log.Config
	Kube  *kube.Config
	Vault *vault.Config
}

func New() *Global {
	return &Global{
		Global: &flags.Global{},
		Env:    nil,
		Log:    nil,
		Kube:   nil,
		Vault:  nil,
	}
}

func (g *Global) Init() error {
	envs, err := InitEnvs()
	if err != nil {
		return err
	}
	g.Env = envs
	g.Log = &log.Config{}
	g.Vault = &vault.Config{
		Role: DefaultVaultRole,
	}
	role := os.Getenv(EnvKeyVaultRole)
	if role != "" {
		g.Vault.Role = role
	}

	if g.Env.RunInPod {
		g.Log.Output = "/tmp"
		g.Vault.Address = DefaultVaultAddr
	} else {
		g.Log.Output = fmt.Sprintf("%s/log/scripts/renew", g.Env.K8SHome)

		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		if hostname == "" {
			return errors.New("get hostname")
		}
	}

	g.Kube = &kube.Config{Kubeconfig: g.KubeConfig}

	return nil
}

func (g *Global) Print() {
	globalv := reflect.ValueOf(*g.Global)
	globalt := reflect.TypeOf(*g.Global)
	log.Info("Global: ")
	for num := 0; num < globalv.NumField(); num++ {
		log.Infof("  %s: %v", globalt.Field(num).Name, globalv.Field(num))
	}
	envsv := reflect.ValueOf(*g.Env)
	envst := reflect.TypeOf(*g.Env)
	log.Info("Envs: ")
	for num := 0; num < envsv.NumField(); num++ {
		log.Infof("  %s: %v", envst.Field(num).Name, envsv.Field(num))
	}
	logv := reflect.ValueOf(*g.Log)
	logt := reflect.TypeOf(*g.Log)
	log.Info("Log: ")
	for num := 0; num < logv.NumField(); num++ {
		log.Infof("  %s: %v", logt.Field(num).Name, logv.Field(num))
	}
	kubev := reflect.ValueOf(*g.Kube)
	kubet := reflect.TypeOf(*g.Kube)
	log.Info("Kube: ")
	for num := 0; num < kubev.NumField(); num++ {
		log.Infof("  %s: %v", kubet.Field(num).Name, kubev.Field(num))
	}
	vaultv := reflect.ValueOf(*g.Vault)
	vaultt := reflect.TypeOf(*g.Vault)
	log.Info("Kube: ")
	for num := 0; num < vaultv.NumField(); num++ {
		log.Infof("  %s: %v", vaultt.Field(num).Name, vaultv.Field(num))
	}
}

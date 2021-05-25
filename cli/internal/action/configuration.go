package action

import (
	"fmt"
	"os"
	"reflect"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

const (
	DefaultVaultAddr = "https://127.0.0.1:8200"
)

type Options struct {
	CertType      string
	Username      string
	Password      string
	SSHKey        string
	Cert          string
	Key           string
	CACert        string
	CDFNamespace  string
	Namespace     string
	Unit          string
	KubeConfig    string
	CAKey         string
	NodeType      string
	Host          string
	OutputDir     string
	ServerCertSan string
	AutoConfirm   bool
	Period        int
}

type Configuration struct {
	*Options

	Env   *env.Settings
	Log   *log.Config
	Kube  *kube.Config
	Vault *vault.Config
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Options: &Options{},
		Env:     nil,
		Log:     nil,
		Kube:    nil,
		Vault:   nil,
	}
}

func (g *Configuration) Init() error {
	envs, err := env.New()
	if err != nil {
		return err
	}
	g.Env = envs
	g.Log = &log.Config{}
	g.Vault = &vault.Config{
		Role: envs.VaultRole,
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

	g.Kube = &kube.Config{Kubeconfig: g.Options.KubeConfig}

	return nil
}

func (g *Configuration) printWithLevel(level string) {
	var print log.LoggerFunc
	if level == log.DefaultLogLevel {
		print = log.Debugf
	} else {
		print = log.Infof
	}
	globalv := reflect.ValueOf(*g.Options)
	globalt := reflect.TypeOf(*g.Options)
	print("Options: ")
	for num := 0; num < globalv.NumField(); num++ {
		print("  %s: %v", globalt.Field(num).Name, globalv.Field(num))
	}
	envsv := reflect.ValueOf(*g.Env)
	envst := reflect.TypeOf(*g.Env)
	print("Envs: ")
	for num := 0; num < envsv.NumField(); num++ {
		print("  %s: %v", envst.Field(num).Name, envsv.Field(num))
	}
	logv := reflect.ValueOf(*g.Log)
	logt := reflect.TypeOf(*g.Log)
	print("Log: ")
	for num := 0; num < logv.NumField(); num++ {
		print("  %s: %v", logt.Field(num).Name, logv.Field(num))
	}
	kubev := reflect.ValueOf(*g.Kube)
	kubet := reflect.TypeOf(*g.Kube)
	print("Kube: ")
	for num := 0; num < kubev.NumField(); num++ {
		print("  %s: %v", kubet.Field(num).Name, kubev.Field(num))
	}
	vaultv := reflect.ValueOf(*g.Vault)
	vaultt := reflect.TypeOf(*g.Vault)
	print("Kube: ")
	for num := 0; num < vaultv.NumField(); num++ {
		print("  %s: %v", vaultt.Field(num).Name, vaultv.Field(num))
	}
}

func (g *Configuration) Print() {
	g.printWithLevel("")
}

func (g *Configuration) Debug() {
	g.printWithLevel(log.DefaultLogLevel)
}

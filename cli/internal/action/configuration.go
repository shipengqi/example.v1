package action

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/pkg/errors"
	"github.com/shipengqi/example.v1/cli/internal/types"

	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

const (
	DefaultVaultAddr     = "https://127.0.0.1:8200"
	DefaultIngressSecret = "nginx-default-secret"
)

type Options struct {
	CertType      string
	Namespace     string
	Username      string
	Password      string
	SSHKey        string
	Cert          string
	Key           string
	CACert        string
	CAKey         string
	Unit          string
	NodeType      string
	Host          string
	OutputDir     string
	ServerCertSan string
	Secret        string
	SkipConfirm   bool
	Local         bool
	Remote        bool
	Validity      int
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
		Options: &Options{
			CertType: types.CertTypeInternal,
		},
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
		g.Vault.Address = fmt.Sprintf("https://%s:8200", hostname)
	}

	g.Kube = &kube.Config{}
	g.CACert = path.Join(g.Env.SSLPath, "ca.crt")
	g.CAKey = path.Join(g.Env.SSLPath, "ca.key")
	g.OutputDir = path.Join(g.Env.SSLPath, "new-certs")
	g.Secret = DefaultIngressSecret

	return nil
}

func (g *Configuration) printWithLevel(level string) {
	var printf log.LoggerFunc
	if level == log.DefaultLogLevel {
		printf = log.Debugf
	} else {
		printf = log.Infof
	}
	globalv := reflect.ValueOf(*g.Options)
	globalt := reflect.TypeOf(*g.Options)
	printf("Options: ")
	for num := 0; num < globalv.NumField(); num++ {
		printf("  %s: %v", globalt.Field(num).Name, globalv.Field(num))
	}
	envsv := reflect.ValueOf(*g.Env)
	envst := reflect.TypeOf(*g.Env)
	printf("Envs: ")
	for num := 0; num < envsv.NumField(); num++ {
		printf("  %s: %v", envst.Field(num).Name, envsv.Field(num))
	}
	logv := reflect.ValueOf(*g.Log)
	logt := reflect.TypeOf(*g.Log)
	printf("Log: ")
	for num := 0; num < logv.NumField(); num++ {
		printf("  %s: %v", logt.Field(num).Name, logv.Field(num))
	}
	kubev := reflect.ValueOf(*g.Kube)
	kubet := reflect.TypeOf(*g.Kube)
	printf("Kube: ")
	for num := 0; num < kubev.NumField(); num++ {
		printf("  %s: %v", kubet.Field(num).Name, kubev.Field(num))
	}
	vaultv := reflect.ValueOf(*g.Vault)
	vaultt := reflect.TypeOf(*g.Vault)
	printf("Vault: ")
	for num := 0; num < vaultv.NumField(); num++ {
		printf("  %s: %v", vaultt.Field(num).Name, vaultv.Field(num))
	}
}

func (g *Configuration) Print() {
	g.printWithLevel("")
}

func (g *Configuration) Debug() {
	g.printWithLevel(log.DefaultLogLevel)
}
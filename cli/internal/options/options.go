package options

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	DefaultMinimumPeriod           = 1
	DefaultRecommendPeriod         = 365
)

type Interface interface {
	Check() error
}

// Global contains all flags for compatibility with older versions
type Global struct {
	CertType      string
	Username      string
	Password      string
	SSHKey        string
	Cert          string
	Key           string
	CACert        string
	CDFNamespace  string
	Namespace     string
	LogLevel      string
	LogOutput     string
	Unit          string
	KubeConfig    string
	CAKey         string
	NodeType      string
	Host          string
	OutputDir     string
	ServerCertSan string
	Install       bool
	Apply         bool
	Renew         bool
	SkipConfirm   bool
	Remote        bool
	Local         bool
	Period        int
}

func (g *Global) Init(flagSet *pflag.FlagSet) {
	flagSet.BoolVarP(
		&g.SkipConfirm,
		"yes",
		"y",
		false,
		"Answer yes for any confirmations.",
	)
	flagSet.StringVarP(
		&g.CertType,
		"type",
		"t",
		"internal",
		"Specifies the type (internal/external) of the server certificates.",
	)
	flagSet.StringVarP(
		&g.Password,
		"password",
		"p",
		"",
		"VM password",
	)
	flagSet.StringVarP(
		&g.Username,
		"username",
		"u",
		"root",
		"VM user",
	)
	flagSet.StringVar(
		&g.SSHKey,
		"key",
		"",
		"SSH key file path.",
	)
	flagSet.IntVarP(
		&g.Period,
		"validity",
		"V", 365,
		"Specifies the validity period (days) of server certificate.",
	)
	flagSet.BoolVar(
		&g.Apply,
		"apply",
		false,
		"Apply certificates.",
	)
	flagSet.BoolVar(
		&g.Renew,
		"renew",
		false,
		"Renew certificates.",
	)
	flagSet.StringVar(
		&g.Cert,
		"tls-cert",
		"",
		"Certificate file path.",
	)
	flagSet.StringVar(
		&g.Key,
		"tls-key",
		"",
		"Private key file path.",
	)
	flagSet.StringVar(
		&g.CACert,
		"tls-cacert",
		"",
		"CA certificate file path.",
	)
	flagSet.StringVar(
		&g.CAKey,
		"tls-cakey",
		"",
		"CA key file path.",
	)
	flagSet.StringVar(
		&g.NodeType,
		"node-type",
		"",
		"Node type (controlplane/worker) of the host which certificates are generated for.",
	)
	flagSet.StringVarP(
		&g.OutputDir,
		"output-dir",
		"d",
		"",
		"The output directory of certificates.",
	)
	flagSet.StringVar(
		&g.Host,
		"host",
		"",
		"The host FQDN or IP address.",
	)
	flagSet.StringVarP(
		&g.Namespace,
		"namespace",
		"n",
		"",
		"Specifies the namespace.",
	)
	flagSet.StringVar(
		&g.CDFNamespace,
		"cdf-namespace",
		"",
		"Specifies the CDF service namespace.",
	)
	flagSet.BoolVar(
		&g.Local,
		"local",
		false,
		"Renew local internal certificates.",
	)
	flagSet.BoolVar(
		&g.Remote,
		"remote",
		false,
		"do not use, just for auto apply certificates.",
	)
	flagSet.BoolVar(
		&g.Install,
		"install",
		false,
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&g.ServerCertSan,
		"server-cert-san",
		"",
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&g.Unit,
		"unit-time",
		"d",
		"unit of time (d/m), just for testing the certificates.",
	)
	flagSet.StringVar(
		&g.KubeConfig,
		"kubeconfig",
		"",
		"Specifies kube config file.",
	)
}

func (g *Global) Check() error {
	if g.CertType != "internal" && g.CertType != "external" {
		return errors.Errorf("The certificates type: %s is invalid.", g.CertType)
	}

	if g.CertType == "external" && g.Local {
		g.Local = false
	}

	if g.NodeType != "" && g.NodeType != types.NodeTypeControlPlane && g.NodeType != types.NodeTypeWorker {
		if g.Host == "" {
			return errors.New("--host flag is missing")
		}
		return errors.Errorf("The node type: %s is invalid.", g.NodeType)
	}

	if g.Period < DefaultMinimumPeriod {
		return errors.Errorf("The minimum period is %d.", DefaultMinimumPeriod)
	}

	if g.Period > DefaultRecommendPeriod {
		log.Warnf("The recommended certificate validity period is or less than %d.", DefaultRecommendPeriod)
	}

	if g.KubeConfig != "" {
		err := os.Setenv("KUBE_CONFIG_FILE", g.KubeConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Global) MarkDeprecated() {
	if g.Renew {
		log.Warn("The '--renew' flag will be deprecated in a future version.")
	}
	if g.Apply {
		log.Warn("The '--apply' flag will be deprecated in a future version.")
	}
	if g.Install {
		log.Warn("The '--install' flag will be deprecated in a future version.")
	}
}

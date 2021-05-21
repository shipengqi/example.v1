package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/cmd/apply"
	"github.com/shipengqi/example.v1/cli/cmd/check"
	configcmd "github.com/shipengqi/example.v1/cli/cmd/config"
	"github.com/shipengqi/example.v1/cli/cmd/create"
	"github.com/shipengqi/example.v1/cli/cmd/renew"
	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

var exitCode = ExitCodeOk
var filename string

func New() *cobra.Command {
	cfg := config.New()

	c := &cobra.Command{
		Use:   "cert-manager",
		Short: "Manages TLS certificates in kubernetes clusters.",
		Long: "To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol.\n" +
			"We uses internal certificates and external certificates to secure its deployment.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := cfg.Init()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Failed to init configuration! ERR:", err)
				os.Exit(ExitCodeError)
			}
			filename, err = log.Init(cfg.Log)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Failed to init logger! ERR:", err)
				os.Exit(ExitCodeError)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !cfg.Remote {
				log.Warn("Additional logging details can be found in:")
				log.Warnf("    %s", filename)
			}
			os.Exit(exitCode)
		},
		Run: func(cmd *cobra.Command, args []string) {
			defer recovery()

			var c action.Interface
			var err error

			if cfg.Renew {
				log.Warn("The '--renew' flag will be deprecated in a future version.")
				if cfg.Install {
					log.Warn("The '--install' flag will be deprecated in a future version.")
					c = action.NewCreate(cfg)
				} else {
					c = action.NewRenew(cfg)
				}
			} else if cfg.Apply {
				log.Warn("The '--renew' flag will be deprecated in a future version.")
				c = action.NewApply(cfg)
			} else {
				log.Info("no matched action flags")
				return
			}

			err = c.Run()
			if err != nil {
				if err == terminal.InterruptErr {
					log.Warnf("%s, interrupted", c.Name())
					return
				}
				exitCode = ExitCodeError
				log.Errorf("%s, ERR: %v", c.Name(), err)
			}
		},
	}

	// Add sub commands
	c.AddCommand(
		create.NewCommand(cfg),
		renew.NewCommand(cfg),
		apply.NewCommand(cfg),
		check.NewCommand(cfg),
		configcmd.NewCommand(cfg),
	)

	cobra.EnableCommandSorting = false
	initFlags(c.Flags(), cfg)

	return c
}

func initFlags(flagSet *pflag.FlagSet, cfg *config.Global) {
	flagSet.BoolVarP(
		&cfg.SkipConfirm,
		"yes",
		"y",
		false,
		"Answer yes for any confirmations.",
	)
	flagSet.StringVarP(
		&cfg.CertType,
		"type",
		"t",
		"internal",
		"Specifies the type (internal/external) of the server certificates.",
	)
	flagSet.StringVarP(
		&cfg.Password,
		"password",
		"p",
		"",
		"VM password",
	)
	flagSet.StringVarP(
		&cfg.Username,
		"username",
		"u",
		"root",
		"VM user",
	)
	flagSet.StringVar(
		&cfg.SSHKey,
		"key",
		"",
		"SSH key file path.",
	)
	flagSet.IntVarP(
		&cfg.Period,
		"validity",
		"V", 365,
		"Specifies the validity period (days) of server certificate.",
	)
	flagSet.BoolVar(
		&cfg.Apply,
		"apply",
		false,
		"Apply certificates.",
	)
	flagSet.BoolVar(
		&cfg.Renew,
		"renew",
		false,
		"Renew certificates.",
	)
	flagSet.StringVar(
		&cfg.Cert,
		"tls-cert",
		"",
		"Certificate file path.",
	)
	flagSet.StringVar(
		&cfg.Key,
		"tls-key",
		"",
		"Private key file path.",
	)
	flagSet.StringVar(
		&cfg.CACert,
		"tls-cacert",
		"",
		"CA certificate file path.",
	)
	flagSet.StringVar(
		&cfg.CAKey,
		"tls-cakey",
		"",
		"CA key file path.",
	)
	flagSet.StringVar(
		&cfg.NodeType,
		"node-type",
		"",
		"Node type (controlplane/worker) of the host which certificates are generated for.",
	)
	flagSet.StringVarP(
		&cfg.OutputDir,
		"output-dir",
		"d",
		"",
		"The output directory of certificates.",
	)
	flagSet.StringVar(
		&cfg.Host,
		"host",
		"",
		"The host FQDN or IP address.",
	)
	flagSet.StringVarP(
		&cfg.Namespace,
		"namespace",
		"n",
		"",
		"Specifies the namespace.",
	)
	flagSet.StringVar(
		&cfg.CDFNamespace,
		"cdf-namespace",
		"",
		"Specifies the CDF service namespace.",
	)
	flagSet.BoolVar(
		&cfg.Local,
		"local",
		false,
		"Renew local internal certificates.",
	)
	flagSet.BoolVar(
		&cfg.Remote,
		"remote",
		false,
		"do not use, just for auto apply certificates.",
	)
	flagSet.BoolVar(
		&cfg.Install,
		"install",
		false,
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&cfg.ServerCertSan,
		"server-cert-san",
		"",
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&cfg.Unit,
		"unit-time",
		"d",
		"unit of time (d/m), just for testing the certificates.",
	)
	flagSet.StringVar(
		&cfg.KubeConfig,
		"kubeconfig",
		"",
		"Specifies kube config file.",
	)
}

func recovery() {
	if err := recover(); err != nil {
		exitCode = ExitCodeError
		log.Errorf("[Recovery] panic recovered: %+v", err)
	}
}

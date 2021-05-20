package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/cmd/apply"
	"github.com/shipengqi/example.v1/cli/cmd/check"
	"github.com/shipengqi/example.v1/cli/cmd/create"
	"github.com/shipengqi/example.v1/cli/cmd/renew"
	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/internal/flags"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

var exitCode = ExitCodeOk
var logFile string

func New() *cobra.Command {
	f := &flags.Global{}
	cfg := &config.Global{
		Flags:        f,
	}
	c := &cobra.Command{
		Use:   "cert-manager",
		Short: "Manages TLS certificates in kubernetes clusters.",
		Long: "To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol.\n" +
			"We uses internal certificates and external certificates to secure its deployment.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if !f.Remote {
				log.Warn("Additional logging details can be found in:")
				log.Warnf("    %s", logFile)
			}
			os.Exit(exitCode)
		},
		Run: func(cmd *cobra.Command, args []string) {
			defer recovery()

			var c action.Interface
			var err error

			if f.Renew {
				log.Warn("The '--renew' flag will be deprecated in a future version.")
				if f.Install {
					log.Warn("The '--install' flag will be deprecated in a future version.")
					c = action.NewCreate(cfg)
				} else {
					c = action.NewRenew(cfg)
				}
			} else if f.Apply {
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
	)

	cobra.EnableCommandSorting = false
	initFlags(c.Flags(), f)

	return c
}

func initFlags(flagSet *pflag.FlagSet, f *flags.Global) {
	flagSet.BoolVarP(
		&f.SkipConfirm,
		"yes",
		"y",
		false,
		"Answer yes for any confirmations.",
	)
	flagSet.StringVarP(
		&f.CertType,
		"type",
		"t",
		"internal",
		"Specifies the type (internal/external) of the server certificates.",
	)
	flagSet.StringVarP(
		&f.Password,
		"password",
		"p",
		"",
		"VM password",
	)
	flagSet.StringVarP(
		&f.Username,
		"username",
		"u",
		"root",
		"VM user",
	)
	flagSet.StringVar(
		&f.SSHKey,
		"key",
		"",
		"SSH key file path.",
	)
	flagSet.IntVarP(
		&f.Period,
		"validity",
		"V", 365,
		"Specifies the validity period (days) of server certificate.",
	)
	flagSet.BoolVar(
		&f.Apply,
		"apply",
		false,
		"Apply certificates.",
	)
	flagSet.BoolVar(
		&f.Renew,
		"renew",
		false,
		"Renew certificates.",
	)
	flagSet.StringVar(
		&f.Cert,
		"tls-cert",
		"",
		"Certificate file path.",
	)
	flagSet.StringVar(
		&f.Key,
		"tls-key",
		"",
		"Private key file path.",
	)
	flagSet.StringVar(
		&f.CACert,
		"tls-cacert",
		"",
		"CA certificate file path.",
	)
	flagSet.StringVar(
		&f.CAKey,
		"tls-cakey",
		"",
		"CA key file path.",
	)
	flagSet.StringVar(
		&f.NodeType,
		"node-type",
		"",
		"Node type (controlplane/worker) of the host which certificates are generated for.",
	)
	flagSet.StringVarP(
		&f.OutputDir,
		"output-dir",
		"d",
		"",
		"The output directory of certificates.",
	)
	flagSet.StringVar(
		&f.Host,
		"host",
		"",
		"The host FQDN or IP address.",
	)
	flagSet.StringVarP(
		&f.Namespace,
		"namespace",
		"n",
		"",
		"Specifies the namespace.",
	)
	flagSet.StringVar(
		&f.CDFNamespace,
		"cdf-namespace",
		"",
		"Specifies the CDF service namespace.",
	)
	flagSet.BoolVar(
		&f.Local,
		"local",
		false,
		"Renew local internal certificates.",
	)
	flagSet.BoolVar(
		&f.Remote,
		"remote",
		false,
		"do not use, just for auto apply certificates.",
	)
	flagSet.BoolVar(
		&f.Install,
		"install",
		false,
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&f.ServerCertSan,
		"server-cert-san",
		"",
		"Just for installing first master node.",
	)
	flagSet.StringVar(
		&f.Unit,
		"unit-time",
		"d",
		"unit of time (d/m), just for testing the certificates.",
	)
	flagSet.StringVar(
		&f.KubeConfig,
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

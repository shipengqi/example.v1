package certmng

import (
	"bytes"
	action2 "github.com/shipengqi/example.v1/apps/cli/internal/action"
	"github.com/shipengqi/example.v1/apps/cli/internal/types"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// HelpTemplate is the help template for cert-manager commands
// This uses the short and long options.
// command should not use this.
// const helpTemplate = `{{.Short}}
// Description:
//   {{.Long}}
// {{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// UsageTemplate is the usage template for cert-manager commands
// This blocks the displaying of the global options. The main cert-manager
// command should not use this.
const usageTemplate = `Usage:{{if (and .Runnable (not .HasAvailableSubCommands))}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.UseLine}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasExample}}

Examples:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var examplesTemplate = `
  SubCommands Mode:
  ./{{.}} renew -t internal -V 365         Renew the internal certificates.
  ./{{.}} renew -t external -V 365         Renew the external certificates.
  ./{{.}} create -t internal -V 365        Create the internal certificates.
  ./{{.}} apply                            Apply the certificates.

  Flags Mode (To be compatible with older versions, will be deprecated in a future version.):
  ./{{.}} --renew -t internal -V 365       Renew the internal certificates.
  ./{{.}} --renew -t external -V 365       Renew the external certificates.
  ./{{.}} --apply                          Apply the certificates.`

const (
	caCertFlagName        = "tls-cacert"
	caKeyFlagName         = "tls-cakey"
	nodeTypeFlagName      = "node-type"
	hostFlagName          = "host"
	serverCertSanFlagName = "server-cert-san"
	certFlagName          = "tls-cert"
	keyFlagName           = "tls-key"
	remoteFlagName        = "remote"
	installFlagName       = "install"
	unitFlagName          = "unit-time"
	renewFlagName         = "renew"
	applyFlagName         = "apply"
	sshKeyFlagName        = "key"
	namespaceFlagName     = "namespace"
	confirmFlagName       = "yes"
	usernameFlagName      = "username"
	passwordFlagName      = "password"
	validityFlagName      = "validity"
	outputFlagName        = "output-dir"
	typeFlagName          = "type"
	cdfnsFlagName         = "cdf-namespace"
	localFlagName         = "local"
	primaryFlagName       = "primary"
	kubeconfigFlagName    = "kubeconfig"
	sourceFlagName        = "resource"
	sourceFieldFlagName   = "field"
)

const (
	rootDesc = `To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol.
We uses internal certificates and external certificates to secure its deployment.`
	passwordFlagDesc = `VM password`
	nodeTypeFlagDesc = "Node type (controlplane/worker) of the host which certificates are generated for."
	typeFlagDesc     = "Specifies the type (internal/external) of the server certificates."
	sourceFlagDesc   = "Specifies the resource type (cm/secret), name. Format: <type>.<name>. e.g. 'cm.tls-cas'"
	validityFlagDesc = "Specifies the validity period (days) of server certificate."
)

type rootOptions struct {
	*renewOptions

	host          string
	serverCertSan string
	remote        bool
	install       bool
	apply         bool
	renew         bool
}

func (o *rootOptions) combine(f *pflag.FlagSet, cfg *action2.Configuration) {
	o.renewOptions.combine(f, cfg)

	if f.Changed(hostFlagName) {
		cfg.Host = o.host
	}
	if f.Changed(serverCertSanFlagName) {
		cfg.ServerCertSan = o.serverCertSan
	}

}

func New(cfg *action2.Configuration) *cobra.Command {
	o := &rootOptions{
		renewOptions: &renewOptions{},
	}

	var buf bytes.Buffer
	baseName := filepath.Base(os.Args[0])
	tmpl, _ := template.New(baseName).Parse(examplesTemplate)

	_ = tmpl.Execute(&buf, baseName)

	c := &cobra.Command{
		Use:     baseName + " [options]",
		Short:   "Manages TLS certificates in CDF clusters.",
		Long:    rootDesc,
		Example: buf.String(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Debugf("Command: %s", strings.Join(os.Args, " "))
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "help" {
				return
			}
			log.Warn("Additional logging details can be found in:")
			log.Warnf("    %s", log.LogFileName)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var c action2.Interface

			if o.renew {
				if o.install {
					c = action2.NewCreate(cfg)
				} else {
					c = action2.NewRenew(cfg)
				}
			} else if o.apply {
				c = action2.NewApply(cfg)
			} else {
				log.Warn("No matched action flags")
				return nil
			}

			log.Debugf("Matched action: %s", c.Name())
			// Cannot call the c.Execute method, Because in go, parent cannot call the child method.
			return action2.Execute(c)
		},
	}

	// Add sub commands
	c.AddCommand(
		newRenewCmd(cfg),
		newCreateCmd(cfg),
		newApplyCmd(cfg),
		newCheckCmd(cfg),
	)

	cobra.EnableCommandSorting = false
	c.Flags().SortFlags = false

	c.SilenceUsage = true
	c.SilenceErrors = true

	addRootFlags(c.Flags(), o)

	c.SetUsageTemplate(usageTemplate)
	// Use:     baseName + " [options]",
	// [options] has been added here, by default, cobra will add [flag] after useline
	c.DisableFlagsInUseLine = true
	return c
}

func addRootFlags(f *pflag.FlagSet, o *rootOptions) {
	f.BoolVarP(&o.skipConfirm, confirmFlagName, "y", false, "Answer yes for any confirmations.")
	f.StringVarP(&o.certType, typeFlagName, "t", "internal", typeFlagDesc)
	f.StringVarP(&o.username, usernameFlagName, "u", "root", "VM user")
	f.StringVarP(&o.password, passwordFlagName, "p", "", passwordFlagDesc)
	f.StringVar(&o.sshKey, sshKeyFlagName, "", "SSH key file path.")
	f.IntVarP(&o.validity, validityFlagName, "V", 365, validityFlagDesc)
	f.BoolVar(&o.apply, applyFlagName, false, "Apply certificates.")
	f.BoolVar(&o.renew, renewFlagName, false, "Renew certificates.")
	f.StringVar(&o.cert, certFlagName, "", "Certificate file path.")
	f.StringVar(&o.key, keyFlagName, "", "Private key file path.")
	f.StringVar(&o.caCert, caCertFlagName, "", "CA certificate file path.")
	f.StringVar(&o.caKey, caKeyFlagName, "", "CA key file path.")
	f.StringVar(&o.nodeType, nodeTypeFlagName, types.NodeTypeControlPlane, nodeTypeFlagDesc)
	f.StringVarP(&o.outputDir, outputFlagName, "d", "", "The output directory of certificates.")
	f.StringVar(&o.host, hostFlagName, "", "The host FQDN or IP address.")
	f.StringVarP(&o.namespace, namespaceFlagName, "n", "", "Specifies the namespace.")
	f.StringVar(&o.cdfNamespace, cdfnsFlagName, "", "Specifies the CDF service namespace.")
	f.BoolVar(&o.local, localFlagName, false, "Renew local internal certificates.")
	f.BoolVar(&o.remote, remoteFlagName, false, "Apply certificates in ssh mode.")
	f.BoolVar(&o.install, installFlagName, false, "Install first master node.")
	f.BoolVar(&o.primary, primaryFlagName, false, "Primary deployment.")
	f.StringVar(&o.serverCertSan, serverCertSanFlagName, "", "server-cert-san for installing first master node.")
	f.StringVar(&o.unit, unitFlagName, "d", "unit of time (d/m), For testing.")
	f.StringVar(&o.kubeconfig, kubeconfigFlagName, "", "Specifies kube config file.")
	f.StringVar(&o.resource, sourceFlagName, action2.SecretNameNginxDefault, "Specifies the resource name(s). Format: <name>,<name>. e.g. '--resource secret1,secret2'")
	f.StringVar(&o.resourceField, sourceFieldFlagName, action2.DefaultResourceKeyTls,
		"Specifies the certificate field of the source data.")

	_ = f.MarkHidden(remoteFlagName)
	_ = f.MarkHidden(installFlagName)
	_ = f.MarkHidden(serverCertSanFlagName)
	_ = f.MarkHidden(unitFlagName)
	_ = f.MarkHidden(cdfnsFlagName)
	_ = f.MarkHidden(primaryFlagName)
	_ = f.MarkHidden(sourceFlagName)
	_ = f.MarkHidden(sourceFieldFlagName)

	_ = f.MarkDeprecated(renewFlagName, "use the 'renew' subcommand instead")
	_ = f.MarkDeprecated(applyFlagName, "use the 'apply' subcommand instead")
	_ = f.MarkDeprecated(installFlagName, "use the 'create' subcommand instead")
}

package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

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
	kubeconfigFlagName    = "kubeconfig"
	sourceFlagName        = "resource"
	sourceFieldFlagName   = "field"
)

const (
	rootDesc = `To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol.
We uses internal certificates and external certificates to secure its deployment.`
	passwordFlagDesc = `VM password (By providing the password with this option, you are disabling or bypassing 
security features, thereby exposing the system to increased security risks. By using this option, 
you understand and agree to assume all associated risks and hold Micro Focus harmless for the same.)`
	nodeTypeFlagDesc = "Node type (controlplane/worker) of the host which certificates are generated for."
	typeFlagDesc     = "Specifies the type (internal/external) of the server certificates."
	sourceFlagDesc   = "Specifies the resource type (cm/secret), name. Format: <type>.<name>. e.g. 'cm.tls-cas'"
	validityFlagDesc = "Specifies the validity period (days) of server certificate."
	examplesDesc     = `
  SubCommands Mode:
  ./cert-manager renew -t internal -V 365         Renew the internal certificates.
  ./cert-manager renew -t external -V 365         Renew the external certificates.
  ./cert-manager create -t internal -V 365        Create the internal certificates.
  ./cert-manager apply                            Apply the certificates.

  Flags Mode (To be compatible with older versions, will be deprecated in a future version.):
  ./cert-manager --renew -t internal -V 365       Renew the internal certificates.
  ./cert-manager --renew -t external -V 365       Renew the external certificates.
  ./cert-manager --apply                          Apply the certificates.`
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

func (o *rootOptions) combine(f *pflag.FlagSet, cfg *action.Configuration) {
	o.renewOptions.combine(f, cfg)

	if f.Changed(hostFlagName) {
		cfg.Host = o.host
	}
	if f.Changed(serverCertSanFlagName) {
		cfg.ServerCertSan = o.serverCertSan
	}

}

func New(cfg *action.Configuration) *cobra.Command {
	o := &rootOptions{
		renewOptions: &renewOptions{},
	}

	c := &cobra.Command{
		Use:     "cert-manager",
		Short:   "Manages TLS certificates in CDF clusters.",
		Long:    rootDesc,
		Example: examplesDesc,
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
			var c action.Interface

			if o.renew {
				if o.install {
					c = action.NewCreate(cfg)
				} else {
					c = action.NewRenew(cfg)
				}
			} else if o.apply {
				c = action.NewApply(cfg)
			} else {
				log.Warn("No matched action flags")
				return nil
			}

			log.Debugf("Matched action: %s", c.Name())
			return action.Execute(c)
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
	f.StringVar(&o.serverCertSan, serverCertSanFlagName, "", "server-cert-san for installing first master node.")
	f.StringVar(&o.unit, unitFlagName, "d", "unit of time (d/m), For testing.")
	f.StringVar(&o.kubeconfig, kubeconfigFlagName, "", "Specifies kube config file.")

	_ = f.MarkHidden(remoteFlagName)
	_ = f.MarkHidden(installFlagName)
	_ = f.MarkHidden(serverCertSanFlagName)
	_ = f.MarkHidden(unitFlagName)

	_ = f.MarkDeprecated(renewFlagName, "Please use the 'renew' subcommand instead")
	_ = f.MarkDeprecated(applyFlagName, "Please use the 'apply' subcommand instead")
	_ = f.MarkDeprecated(installFlagName, "Please use the 'create' subcommand instead")
}

package certmng

import (
	"bytes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	rootDesc = `To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol.
We uses internal certificates and external certificates to secure its deployment.`
	passwordFlagDesc = `VM password (By providing the password with this option, you are disabling or bypassing 
security features, thereby exposing the system to increased security risks. By using this option, 
you understand and agree to assume all associated risks and hold Micro Focus harmless for the same.`
	nodeTypeFlagDesc = "Node type (controlplane/worker) of the host which certificates are generated for."
	typeFlagDesc     = "Specifies the type (internal/external) of the server certificates."
	validityFlagDesc = "Specifies the validity period (days) of server certificate."
	examplesDesc     = `
  SubCommands Mode:
  ./renewCert renew -t internal -V 365         Renew the internal certificates.
  ./renewCert renew -t external -V 365         Renew the external certificates.
  ./renewCert create -t internal -V 365        Create the internal certificates.
  ./renewCert apply                            Apply the certificates.

  Flags Mode (To be compatible with older versions, will be deprecated in a future version.):
  ./renewCert --renew -t internal -V 365       Renew the internal certificates.
  ./renewCert --renew -t external -V 365       Renew the external certificates.
  ./renewCert --apply                          Apply the certificates.`
)

type rootOptions struct {
	certType      string
	username      string
	password      string
	sshKey        string
	cert          string
	key           string
	caCert        string
	cdfNamespace  string
	namespace     string
	unit          string
	kubeconfig    string
	caKey         string
	nodeType      string
	host          string
	outputDir     string
	serverCertSan string
	install       bool
	apply         bool
	renew         bool
	skipConfirm   bool
	remote        bool
	local         bool
	validity      int
}

func New(cfg *action.Configuration) *cobra.Command {
	o := &rootOptions{}

	c := &cobra.Command{
		Use:     "cert-manager",
		Short:   "Manages TLS certificates in kubernetes clusters.",
		Long:    rootDesc,
		Example: examplesDesc,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg.Debug()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			log.Warn("Additional logging details can be found in:")
			log.Warnf("    %s", log.LogFileName)
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
				log.Info("No matched action flags")
				return nil
			}

			log.Infof("Matched action flag: %s", c.Name())
			return c.Execute()
		},
	}

	// Add sub commands
	c.AddCommand(
		newRenewCmd(cfg),
		newCreateCmd(cfg),
		newApplyCmd(cfg),
		newCheckCmd(cfg),
		newConfigCmd(cfg),
	)

	cobra.EnableCommandSorting = false

	addFlags(c.Flags(), o)

	return c
}

func addFlags(f *pflag.FlagSet, o *rootOptions) {
	f.BoolVarP(&o.skipConfirm, "yes", "y", false, "Answer yes for any confirmations.")
	f.StringVarP(&o.certType, "type", "t", "internal", typeFlagDesc)
	f.StringVarP(&o.username, "username", "u", "root", "VM user")
	f.StringVarP(&o.password, "password", "p", "", passwordFlagDesc)
	f.StringVar(&o.sshKey, "key", "", "SSH key file path.")
	f.IntVarP(&o.validity, "validity", "V", 365, validityFlagDesc)
	f.BoolVar(&o.apply, "apply", false, "Apply certificates.")
	f.BoolVar(&o.renew, "renew", false, "Renew certificates.")
	f.StringVar(&o.cert, "tls-cert", "", "Certificate file path.")
	f.StringVar(&o.key, "tls-key", "", "Private key file path.")
	f.StringVar(&o.caCert, "tls-cacert", "", "CA certificate file path.")
	f.StringVar(&o.caKey, "tls-cakey", "", "CA key file path.")
	f.StringVar(&o.nodeType, "node-type", "", nodeTypeFlagDesc)
	f.StringVarP(&o.outputDir, "output-dir", "d", "", "The output directory of certificates.")
	f.StringVar(&o.host, "host", "", "The host FQDN or IP address.")
	f.StringVarP(&o.namespace, "namespace", "n", "", "Specifies the namespace.")
	f.StringVar(&o.cdfNamespace, "cdf-namespace", "", "Specifies the CDF service namespace.")
	f.BoolVar(&o.local, "local", false, "Renew local internal certificates.")
	f.BoolVar(&o.remote, "remote", false, "Apply certificates in ssh mode.")
	f.BoolVar(&o.install, "install", false, "Install first master node.")
	f.StringVar(&o.serverCertSan, "server-cert-san", "", "server-cert-san for installing first master node.")
	f.StringVar(&o.unit, "unit-time", "d", "unit of time (d/m), For testing.")
	f.StringVar(&o.kubeconfig, "kubeconfig", "", "Specifies kube config file.")

	_ = f.MarkHidden("remote")
	_ = f.MarkHidden("install")
	_ = f.MarkHidden("server-cert-san")
	_ = f.MarkHidden("unit-time")
	_ = f.MarkHidden("server-cert-san")

	_ = f.MarkDeprecated("renew", "'renew' flag will be deprecated in a future version.")
	_ = f.MarkDeprecated("apply", "'apply' flag will be deprecated in a future version.")
	_ = f.MarkDeprecated("install", "'install' flag will be deprecated in a future version.")
}

func combine(cmd *cobra.Command, args []string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(cmd.Name())

	for k := range args {
		buf.WriteString(" ")
		buf.WriteString(args[k])
	}
	return buf.String()
}

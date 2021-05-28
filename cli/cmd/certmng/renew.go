package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type renewOptions struct {
	certType     string
	username     string
	password     string
	sshKey       string
	cert         string
	key          string
	caCert       string
	cdfNamespace string
	namespace    string
	unit         string
	kubeconfig   string
	caKey        string
	nodeType     string
	host         string
	outputDir    string
	secret       string
	skipConfirm  bool
	local        bool
	validity     int
}

func (o *renewOptions) combine(f *pflag.FlagSet, cfg *action.Configuration) {
	if f.Changed(confirmFlagName) {
		cfg.SkipConfirm = o.skipConfirm
	}
	if f.Changed(typeFlagName) {
		cfg.CertType = o.certType
	}
	if f.Changed(usernameFlagName) {
		cfg.Username = o.username
	}
	if f.Changed(passwordFlagName) {
		cfg.Password = o.password
	}
	if f.Changed(sshKeyFlagName) {
		cfg.SSHKey = o.sshKey
	}
	if f.Changed(validityFlagName) {
		cfg.Validity = o.validity
	}
	if f.Changed(certFlagName) {
		cfg.Cert = o.cert
	}
	if f.Changed(keyFlagName) {
		cfg.Key = o.key
	}
	if f.Changed(caCertFlagName) {
		cfg.CACert = o.caCert
	}
	if f.Changed(caKeyFlagName) {
		cfg.CAKey = o.caKey
	}
	if f.Changed(nodeTypeFlagName) {
		cfg.NodeType = o.nodeType
	}
	if f.Changed(outputFlagName) {
		cfg.OutputDir = o.outputDir
	}
	if f.Changed(hostFlagName) {
		cfg.Host = o.host
	}
	if f.Changed(namespaceFlagName) {
		cfg.Namespace = o.namespace
	}
	if f.Changed(localFlagName) {
		cfg.Local = o.local
	}
	if f.Changed(unitFlagName) {
		cfg.Unit = o.unit
	}

	if f.Changed(cdfnsFlagName) {
		cfg.Env.CDFNamespace = o.cdfNamespace
	}
	if f.Changed(kubeconfigFlagName) {
		cfg.Kube.Kubeconfig = o.kubeconfig
	}

	if len(cfg.Namespace) == 0 {
		cfg.Namespace = cfg.Env.CDFNamespace
	}
}

func newRenewCmd(cfg *action.Configuration) *cobra.Command {
	o := &renewOptions{}
	c := &cobra.Command{
		Use:   renewFlagName,
		Short: "Renew the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			r := action.NewRenew(cfg)
			return r.Execute()
		},
	}

	addRenewFlags(c.Flags(), o)

	return c
}

func addRenewFlags(f *pflag.FlagSet, o *renewOptions) {
	f.BoolVarP(&o.skipConfirm, confirmFlagName, "y", false, "Answer yes for any confirmations.")
	f.StringVarP(&o.certType, typeFlagName, "t", "internal", typeFlagDesc)
	f.StringVarP(&o.username, usernameFlagName, "u", "root", "VM user")
	f.StringVarP(&o.password, passwordFlagName, "p", "", passwordFlagDesc)
	f.StringVar(&o.sshKey, sshKeyFlagName, "", "SSH key file path.")
	f.IntVarP(&o.validity, validityFlagName, "V", 365, validityFlagDesc)
	f.StringVar(&o.cert, certFlagName, "", "Certificate file path.")
	f.StringVar(&o.key, keyFlagName, "", "Private key file path.")
	f.StringVar(&o.caCert, caCertFlagName, "", "CA certificate file path.")
	f.StringVar(&o.caKey, caKeyFlagName, "", "CA key file path.")
	f.StringVar(&o.nodeType, nodeTypeFlagName, "", nodeTypeFlagDesc)
	f.StringVarP(&o.outputDir, outputFlagName, "d", "", "The output directory of certificates.")
	f.StringVar(&o.host, hostFlagName, "", "The host FQDN or IP address.")
	f.StringVarP(&o.namespace, namespaceFlagName, "n", "", "Specifies the namespace.")
	f.StringVar(&o.cdfNamespace, cdfnsFlagName, "", "Specifies the CDF service namespace.")
	f.BoolVar(&o.local, localFlagName, false, "Renew local internal certificates.")
	f.StringVar(&o.unit, unitFlagName, "d", "unit of time (d/m), For testing.")
	f.StringVar(&o.kubeconfig, kubeconfigFlagName, "", "Specifies kube config file.")

	_ = f.MarkHidden(unitFlagName)
}

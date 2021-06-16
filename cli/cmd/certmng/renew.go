package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/types"
)

type renewOptions struct {
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
	outputDir     string
	resource      string
	resourceField string
	skipConfirm   bool
	local         bool
	primary       bool
	validity      int
}

func (o *renewOptions) combine(f *pflag.FlagSet, cfg *action.Configuration) {
	if f.Changed(passwordFlagName) {
		cfg.Password = o.password
	}
	if f.Changed(sshKeyFlagName) {
		cfg.SSHKey = o.sshKey
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

	if f.Changed(outputFlagName) {
		cfg.OutputDir = o.outputDir
	}
	if f.Changed(namespaceFlagName) {
		cfg.Namespace = o.namespace
	}

	if f.Changed(cdfnsFlagName) {
		cfg.Env.CDFNamespace = o.cdfNamespace
	}
	if f.Changed(kubeconfigFlagName) {
		cfg.Kube.Kubeconfig = o.kubeconfig
	}
	if f.Changed(primaryFlagName) {
		cfg.Cluster.IsPrimary = o.primary
	}
	if len(cfg.Namespace) == 0 {
		cfg.Namespace = cfg.Env.CDFNamespace
	}

	// default value
	cfg.SkipConfirm = o.skipConfirm
	cfg.CertType = o.certType
	cfg.Username = o.username
	cfg.Validity = o.validity
	cfg.NodeType = o.nodeType
	cfg.Local = o.local
	cfg.Unit = o.unit
	cfg.Resource = o.resource
	cfg.ResourceField = o.resourceField
}

func newRenewCmd(cfg *action.Configuration) *cobra.Command {
	o := &renewOptions{}
	c := &cobra.Command{
		Use:   renewFlagName + " [options]",
		Short: "Renew the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			r := action.NewRenew(cfg)
			return action.Execute(r)
		},
	}

	c.Flags().SortFlags = false
	c.DisableFlagsInUseLine = true
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
	f.StringVar(&o.nodeType, nodeTypeFlagName, types.NodeTypeControlPlane, nodeTypeFlagDesc)
	f.StringVarP(&o.outputDir, outputFlagName, "d", "", "The output directory of certificates.")
	f.StringVarP(&o.namespace, namespaceFlagName, "n", "", "Specifies the namespace.")
	f.StringVar(&o.cdfNamespace, cdfnsFlagName, "", "Specifies the CDF service namespace.")
	f.BoolVar(&o.local, localFlagName, false, "Renew local internal certificates.")
	f.BoolVar(&o.primary, primaryFlagName, false, "Primary deployment.")
	f.StringVar(&o.unit, unitFlagName, "d", "unit of time (d/m), For testing.")
	f.StringVar(&o.kubeconfig, kubeconfigFlagName, "", "Specifies kube config file.")
	f.StringVar(&o.resource, sourceFlagName, action.SecretNameNginxDefault, "Specifies the resource name(s). Format: <name>,<name>. e.g. '--resource secret1,secret2'")
	f.StringVar(&o.resourceField, sourceFieldFlagName, action.DefaultResourceKeyTls,
		"Specifies the certificate field of the source data.")

	_ = f.MarkHidden(unitFlagName)
	_ = f.MarkHidden(cdfnsFlagName)
}

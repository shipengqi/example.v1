package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type createOptions struct {
	caCert        string
	caKey         string
	nodeType      string
	host          string
	serverCertSan string
	outputDir     string
}

func (o *createOptions) combine(f *pflag.FlagSet, cfg *action.Configuration)  {
	if f.Changed(caCertFlagName) {
		cfg.CACert = o.caCert
	}
	if f.Changed(caKeyFlagName) {
		cfg.CAKey = o.caKey
	}
	if f.Changed(nodeTypeFlagName) {
		cfg.NodeType = o.nodeType
	}
	if f.Changed(hostFlagName) {
		cfg.Host = o.host
	}
	if f.Changed(serverCertSanFlagName) {
		cfg.ServerCertSan = o.serverCertSan
	}
	if f.Changed(outputFlagName) {
		cfg.OutputDir = o.outputDir
	}
}

func newCreateCmd(cfg *action.Configuration) *cobra.Command {
	o := &createOptions{}
	c := &cobra.Command{
		Use:   "create",
		Short: "Create the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCreate(cfg)
			return c.Execute()
		},
	}

	f := c.Flags()
	f.StringVar(&o.caCert, caCertFlagName, "", "CA certificate file path.")
	f.StringVar(&o.caKey, caKeyFlagName, "", "CA key file path.")
	f.StringVar(&o.nodeType, nodeTypeFlagName, "", nodeTypeFlagDesc)
	f.StringVar(&o.host, hostFlagName, "", "The host FQDN or IP address.")
	f.StringVar(&o.serverCertSan, serverCertSanFlagName, "", "server-cert-san for installing first master node.")

	return c
}
package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type checkOptions struct {
	cert      string
	secret    string
	namespace string
}

func (o *checkOptions) combine(f *pflag.FlagSet, cfg *action.Configuration)  {
	if f.Changed(certFlagName) {
		cfg.Cert = o.cert
	}
	if f.Changed(namespaceFlagName) {
		cfg.Namespace = o.namespace
	}
	if f.Changed(secretFlagName) {
		cfg.Secret = o.secret
	}
}

func newCheckCmd(cfg *action.Configuration) *cobra.Command {
	o := &checkOptions{}
	c := &cobra.Command{
		Use:   "check",
		Short: "Check the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			f := cmd.Flags()
			o.combine(f, cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCheck(cfg)
			return c.Execute()
		},
	}

	f := c.Flags()
	f.StringVar(&o.cert, certFlagName, "", "Certificate file path.")
	f.StringVar(&o.secret, secretFlagName, "", "Specifies the secret name.")
	f.StringVarP(&o.namespace, namespaceFlagName, "n", "", "Specifies the namespace.")

	return c
}

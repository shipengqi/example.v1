package certmng

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type checkOptions struct {
	cert          string
	resource      string
	resourceField string
	namespace     string
}

func (o *checkOptions) combine(f *pflag.FlagSet, cfg *action.Configuration) {
	if f.Changed(certFlagName) {
		cfg.Cert = o.cert
	}
	if f.Changed(namespaceFlagName) {
		cfg.Namespace = o.namespace
	}
	if f.Changed(sourceFlagName) {
		cfg.Resource = o.resource
	}
	if f.Changed(sourceFieldFlagName) {
		cfg.ResourceField = o.resourceField
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
			return action.Execute(c)
		},
	}
	c.Flags().SortFlags = false
	c.DisableFlagsInUseLine = true
	f := c.Flags()
	f.StringVar(&o.cert, certFlagName, "", "Certificate file path.")
	f.StringVar(&o.resource, sourceFlagName, "", sourceFlagDesc)
	f.StringVar(&o.resourceField, sourceFieldFlagName, "",
		"Specifies the certificate field of the source data.")
	f.StringVarP(&o.namespace, namespaceFlagName, "n", "", "Specifies the namespace.")

	return c
}

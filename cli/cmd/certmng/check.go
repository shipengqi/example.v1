package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type checkOptions struct {
	cert      string
	secret    string
	namespace string
}

func newCheckCmd(cfg *action.Configuration) *cobra.Command {
	o := &checkOptions{}
	c := &cobra.Command{
		Use:    "check",
		Short:  "Check the internal/external certificates in CDF clusters.",
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

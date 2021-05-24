package certmng

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/env"
)

type checkOptions struct {
	CertType       string
	Cert           string
	Key            string
	CACert         string
	CDFNamespace   string
	Namespace      string
}

func newCheckCmd(cfg *env.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "check",
		Short: "Check the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCheck(cfg)
			err := c.Run()
			if err != nil {
				return errors.Wrapf(err, "%s.Run()", c.Name())
			}
			return nil
		},
	}

	return c
}

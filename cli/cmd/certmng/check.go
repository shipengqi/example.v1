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
	c := &cobra.Command{
		Use:    "check",
		Short:  "Check the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCheck(cfg)
			return c.Execute()
		},
	}

	return c
}

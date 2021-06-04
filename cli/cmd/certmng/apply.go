package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

var remote bool

func newApplyCmd(cfg *action.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:    applyFlagName,
		Short:  "Apply the internal/external certificates in CDF clusters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewApply(cfg)
			return a.Execute()
		},
	}

	f := c.Flags()
	f.BoolVar(&remote, remoteFlagName, false, "apply certificates in ssh mode")
	_ = f.MarkHidden(remoteFlagName)

	return c
}

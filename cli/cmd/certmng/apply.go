package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

var remote bool

func newApplyCmd(cfg *action.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:    "apply",
		Short:  "Apply the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if remote {
				return
			}
			log.Warn("Additional logging details can be found in:")
			log.Warnf("    %s", log.LogFileName)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewApply(cfg)
			return a.Execute()
		},
	}

	f := c.Flags()
	f.BoolVar(&remote, "remote", false, "apply certificates in ssh mode")
	_ = f.MarkHidden("remote")

	return c
}

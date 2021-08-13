package certmng

import (
	action2 "github.com/shipengqi/example.v1/apps/cli/internal/action"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"github.com/spf13/cobra"
)

var remote bool

func newApplyCmd(cfg *action2.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   applyFlagName + " [options]",
		Short: "Apply the internal/external certificates in CDF clusters.",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cfg.Remote {
				return
			}
			log.Warn("Additional logging details can be found in:")
			log.Warnf("    %s", log.LogFileName)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action2.NewApply(cfg)
			return action2.Execute(a)
		},
	}

	c.DisableFlagsInUseLine = true
	f := c.Flags()
	f.BoolVar(&remote, remoteFlagName, false, "apply certificates in ssh mode")
	_ = f.MarkHidden(remoteFlagName)

	return c
}

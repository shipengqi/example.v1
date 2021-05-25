package certmng

import (
	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/spf13/cobra"
)

func newConfigCmd(cfg *action.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Dump the cert-manager configurations.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
	}
	return c
}

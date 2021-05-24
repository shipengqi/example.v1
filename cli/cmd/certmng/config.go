package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/env"
)

func newConfigCmd(cfg *env.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Dump the cert-manager configurations.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}
	return c
}

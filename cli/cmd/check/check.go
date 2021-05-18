package check

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func NewCommand(cfg *config.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "check",
		Short: "Check the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := action.NewCheck(cfg)
			err := c.Run()
			if err != nil {
				log.Errorf("check, ERR: %v", err)
			}
			return
		},
	}

	return c
}

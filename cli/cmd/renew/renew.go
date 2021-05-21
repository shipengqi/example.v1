package renew

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func NewCommand(cfg *config.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "renew",
		Short: "Renew the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
		Run: func(cmd *cobra.Command, args []string) {
			r := action.NewRenew(cfg)
			err := r.Run()
			if err != nil {
				log.Errorf("renew, ERR: %v", err)
			}
			return
		},
	}

	return c
}

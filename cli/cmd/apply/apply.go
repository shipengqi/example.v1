package apply

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func NewCommand(cfg *config.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "apply",
		Short: "Apply the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
		Run: func(cmd *cobra.Command, args []string) {
			a := action.NewApply(cfg)
			err := a.Run()
			if err != nil {
				log.Errorf("apply, ERR: %v", err)
				log.Warnf("Make sure that you have run the '%s/scripts/renewCert --apply' "+
					"on other master nodes.", cfg.Env.K8SHome)
			}
			return
		},
	}

	return c
}

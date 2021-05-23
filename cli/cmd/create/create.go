package create

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/config"
)

func NewCommand(cfg *config.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Create the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCreate(cfg)
			err := c.Run()
			if err != nil {
				return errors.Wrapf(err, "%s.Run()", c.Name())
			}
			return nil
		},
	}

	return c
}

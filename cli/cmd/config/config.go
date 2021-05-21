package config

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/config"
)

func NewCommand(cfg *config.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Dump the cert-manager configurations.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg.Print()
		},
	}

	return c
}

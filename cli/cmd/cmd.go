package cmd

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/cmd/apply"
	"github.com/shipengqi/example.v1/cli/cmd/check"
	"github.com/shipengqi/example.v1/cli/cmd/create"
	"github.com/shipengqi/example.v1/cli/cmd/renew"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "cert-manager",
		Short: "Manages TLS certificates in kubernetes clusters.",
		Long: "To securely deploy the kubernetes, we recommend that you use the TLS/SSL communication protocol. " +
			"We uses internal certificates and external certificates to secure its deployment.",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Add sub commands
	c.AddCommand(
		create.NewCommand(),
		renew.NewCommand(),
		apply.NewCommand(),
		check.NewCommand(),
	)

	cobra.EnableCommandSorting = false

	return c
}

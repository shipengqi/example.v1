package apply

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "apply",
		Short: "Apply the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return c
}

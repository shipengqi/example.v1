package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type createOptions struct {
	CAKey          string
	NodeType       string
	Host           string
	KubeApiCertSan string
}

func newCreateCmd(cfg *action.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "Create the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := action.NewCreate(cfg)
			return c.Execute()
		},
	}

	return c
}

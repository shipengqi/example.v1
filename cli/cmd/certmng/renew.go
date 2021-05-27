package certmng

import (
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
)

type renewOptions struct {
	CertType      string
	Username      string
	Password      string
	SSHKey        string
	Cert          string
	Key           string
	CACert        string
	CDFNamespace  string
	Namespace     string
	Unit          string
	KubeConfig    string
	CAKey         string
	NodeType      string
	Host          string
	OutputDir     string
	SkipConfirm   bool
	Local         bool
	Period        int
}

func newRenewCmd(cfg *action.Configuration) *cobra.Command {
	c := &cobra.Command{
		Use:   renewFlagName,
		Short: "Renew the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			r := action.NewRenew(cfg)
			return r.Execute()
		},
	}

	return c
}

package certmng

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/env"
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
	LogLevel      string
	LogOutput     string
	Unit          string
	KubeConfig    string
	CAKey         string
	NodeType      string
	Host          string
	OutputDir     string
	ServerCertSan string
	Install       bool
	Apply         bool
	Renew         bool
	SkipConfirm   bool
	Remote        bool
	Local         bool
	Period        int
}

func newRenewCmd(cfg *env.Global) *cobra.Command {
	c := &cobra.Command{
		Use:   "renew",
		Short: "Renew the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			r := action.NewRenew(cfg)
			err := r.Run()
			if err != nil {
				return errors.Wrapf(err, "%s.Run()", r.Name())
			}
			return nil
		},
	}

	return c
}

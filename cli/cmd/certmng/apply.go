package certmng

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/example.v1/cli/internal/action"
	"github.com/shipengqi/example.v1/cli/internal/env"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type applyOptions struct {
	Remote bool
}

func newApplyCmd(cfg *env.Global) *cobra.Command {
	c := &cobra.Command{
		Use:    "apply",
		Short:  "Apply the internal/external certificates in CDF clusters.",
		PreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := action.NewApply(cfg)
			err := a.Run()
			if err != nil {
				log.Warnf("Make sure that you have run the '%s/scripts/renewCert --apply' "+
					"on other master nodes.", cfg.Env.K8SHome)
				return errors.Wrapf(err, "%s.Run()", a.Name())
			}
			return nil
		},
	}

	return c
}

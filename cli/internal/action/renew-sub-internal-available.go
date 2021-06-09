package action

import (
	"strings"

	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubInternalAvailable struct {
	*action

	kube *kube.Client
}

func NewRenewSubInternalAvailable(cfg *Configuration) Interface {
	c := &renewSubInternalAvailable{
		action: &action{
			name: "renew-sub-internal-available",
			cfg:  cfg,
		},
	}

	return c
}

func (a *renewSubInternalAvailable) Name() string {
	return a.name
}

func (a *renewSubInternalAvailable) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))

	return sysc.RenewRERemoteExecution(a.cfg.Env.CDFNamespace, a.cfg.Namespace,
		a.cfg.Unit, a.cfg.Validity, a.cfg.SkipConfirm)
}

func (a *renewSubInternalAvailable) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	cm, err := a.kube.GetConfigMap(a.cfg.Env.CDFNamespace, ConfigMapNameCDFCluster)
	if err != nil {
		log.Warnf("kube.GetConfigMap(): %v", err)
	} else {
		a.cfg.Cluster.VirtualIP = cm.Data["HA_VIRTUAL_IP"]
		a.cfg.Cluster.LoadBalanceIP = cm.Data["LOAD_BALANCER_HOST"]
		a.cfg.Cluster.KubeServiceIP = cm.Data["K8S_DEFAULT_SVC_IP"]
		a.cfg.Cluster.FirstMasterNode = cm.Data["FIRST_MASTER_NODE"]
		a.cfg.Cluster.EtcdEndpoint = cm.Data["ETCD_ENDPOINT"]
	}

	a.cfg.Debug()

	return nil
}

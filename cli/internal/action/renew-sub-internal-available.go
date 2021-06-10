package action

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/node"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type renewSubInternalAvailable struct {
	*action

	generator certs.Generator
}

func NewRenewSubInternalAvailable(cfg *Configuration) Interface {
	c := &renewSubInternalAvailable{
		action: newActionWithKube("renew-sub-internal-available", cfg),
	}

	key, err := c.parseCAKey()
	if err != nil {
		panic(err)
	}

	g, err := infra.New(cfg.CACert, key)
	if err != nil {
		panic(err)
	}
	c.generator = g

	return c
}

func (a *renewSubInternalAvailable) Name() string {
	return a.name
}

func (a *renewSubInternalAvailable) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	nodes, err := a.getNodes()
	if err != nil {
		return err
	}

	if len(nodes) < 1 {
		return errors.New("get node 0")
	}
	if a.cfg.Env.RunOnMaster {
		err := a.iterateSecrets(a.generator)
		if err != nil {
			return err
		}
	}

	return nil
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

func (a *renewSubInternalAvailable) getNodes() ([]node.Node, error) {
	var nodes []node.Node
	ns, err := a.kube.GetNodes()
	if err != nil {
		return []node.Node{}, err
	}
	for _, v := range ns.Items {
		m, ok := v.Labels["master"]
		if ok {
			isMaster, err := strconv.ParseBool(m)
			if err != nil {
				log.Warnf("ParseBool, Err: %s", err)
			}
			nodes = append(nodes, node.Node{Address: v.Name, Master: isMaster})
		} else {
			nodes = append(nodes, node.Node{Address: v.Name, Master: false})
		}
	}
	return nodes, nil
}

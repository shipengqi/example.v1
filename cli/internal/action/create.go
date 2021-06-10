package action

import (
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type create struct {
	*action

	generator certs.Generator
}

func NewCreate(cfg *Configuration) Interface {
	c := &create{
		action: newActionWithKube("create", cfg),
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

func (a *create) Name() string {
	return a.name
}

func (a *create) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	log.Info("Creating certificates ...")
	var master bool

	switch a.cfg.NodeType {
	case types.NodeTypeControlPlane:
		master = true
		break
	case types.NodeTypeWorker:
		master = false
		break
	default:
		return errors.Errorf("unknown node type: %s", a.cfg.NodeType)
	}
	return a.iterate(a.cfg.Host, master, false, a.generator)
}

func (a *create) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	cm, err := a.kube.GetConfigMap(a.cfg.Env.CDFNamespace, ConfigMapNameCDFCluster)
	if err != nil {
		log.Warnf("kube.GetConfigMap(): %v", err)
	} else {
		a.cfg.Cluster.VirtualIP = cm.Data["HA_VIRTUAL_IP"]
		a.cfg.Cluster.LoadBalanceIP = cm.Data["LOAD_BALANCER_HOST"]
		a.cfg.Cluster.KubeServiceIP = cm.Data["K8S_DEFAULT_SVC_IP"]
	}

	a.cfg.Debug()

	// create new-certs folder for internal cert
	return os.MkdirAll(a.cfg.OutputDir, 0744)
}

func (a *create) PostRun() error {
	log.Debugf("***** %s PostRun *****", strings.ToUpper(a.name))
	log.Info("Finished!")
	return nil
}

package action

import (
	"os"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type create struct {
	*action

	kube      *kube.Client
	generator certs.Generator
}

func NewCreate(cfg *Configuration) Interface {
	g, err := infra.New(cfg.CACert, cfg.CAKey)
	if err != nil {
		panic(err)
	}
	kclient, err := kube.New(cfg.Kube)
	if err != nil {
		panic(err)
	}

	c := &create{
		action: &action{
			name: "create",
			cfg:  cfg,
		},
		kube:      kclient,
		generator: g,
	}

	return c
}

func (a *create) Name() string {
	return a.name
}

func (a *create) Run() error {
	log.Debug("*****  CREATE CRT  *****")
	log.Info("Creating certificates ...")
	var isMater bool

	switch a.cfg.NodeType {
	case types.NodeTypeControlPlane:
		isMater = true
		break
	case types.NodeTypeWorker:
		isMater = false
		break
	default:
		return errors.Errorf("unknown node type: %s", a.cfg.NodeType)
	}
	return a.iterate(a.cfg.Host, isMater, a.generator)
}

func (a *create) PreRun() error {
	log.Debug("*****  CREATE PRE RUN  *****")

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
	log.Info("Finished.")
	return nil
}

func (a *create) Execute() error {
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}

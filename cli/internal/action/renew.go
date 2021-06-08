package action

import (
	"fmt"
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/deployment"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var DropError = errors.New("Exit")

type renew struct {
	*action

	expired   bool
	kube      *kube.Client
	generator certs.Generator
}

func NewRenew(cfg *Configuration) Interface {
	var (
		g   certs.Generator
		err error
	)

	kclient, err := kube.New(cfg.Kube)
	if err != nil {
		panic(err)
	}

	if cfg.CertType == types.CertTypeExternal {
		cas, err := getCAs(kclient, cfg.Env.CDFNamespace)
		if err != nil {
			panic(err)
		}
		cfg.Vault.CAs = cas

		g, err = deployment.New(cfg.Env.CDFNamespace, cfg.Kube, cfg.Vault)
	} else {
		g, err = infra.New(cfg.CACert, cfg.CAKey)
	}
	if err != nil {
		panic(err)
	}

	return &renew{
		action: &action{
			name: "renew",
			cfg:  cfg,
		},
		kube:      kclient,
		generator: g,
	}
}

func (a *renew) Name() string {
	return a.name
}

func (a *renew) PreRun() error {
	log.Debug("*****  RENEW PRE RUN  *****")

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

	cm, err = a.kube.GetConfigMap(a.cfg.Env.CDFNamespace, ConfigMapNameCDF)
	if err != nil {
		log.Warnf("kube.GetConfigMap(): %v", err)
	} else {
		a.cfg.Cluster.ExternalHost = cm.Data["EXTERNAL_ACCESS_HOST"]
	}

	a.cfg.Debug()

	// create new-certs folder for internal cert
	err = os.MkdirAll(a.cfg.OutputDir, 0744)
	if err != nil {
		return err
	}

	// check cert validity
	crt, err := utils.ParseCrt(a.cfg.Cert)
	if err != nil {
		return err
	}
	available := utils.CheckCrtValidity(crt)
	if available <= 0 {
		log.Infof("The certificate: %s has already expired.", a.cfg.Cert)
	} else {
		log.Infof("The certificate: %s will expire in %d hour(s).", a.cfg.Cert, available)
	}

	// confirm
	if a.cfg.SkipConfirm {
		return nil
	}
	if a.cfg.Env.RunInPod {
		return nil
	}
	confirm, err := prompt.Confirm("Are you sure to continue")
	if err != nil {
		return err
	}
	if !confirm {
		return DropError
	}
	return nil
}

func (a *renew) Run() error {
	log.Debug("*****  RENEW CRT  *****")
	log.Info("Renewing certificates ...")

	switch a.cfg.CertType {
	case types.CertTypeInternal:
		log.Infof("Cert type: %s", types.CertTypeInternal)
		return a.generator.GenAndDump(&certs.Certificate{
			CN:       a.cfg.Host,
			UintTime: a.cfg.Unit,
			Validity: a.cfg.Validity,
		}, a.cfg.OutputDir)
	case types.CertTypeExternal:
		log.Infof("Cert type: %s", types.CertTypeExternal)
		return a.renewExternal()
	default:
		return errors.Errorf("unknown cert type: %s", a.cfg.CertType)
	}
}

func (a *renew) PostRun() error {
	log.Info("Finished.")
	return nil
}

func (a *renew) Execute() error {
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

func (a *renew) renewExternal() error {
	if len(a.cfg.Cert) > 0 && len(a.cfg.Key) > 0 {
		log.Info("Renewing custom certificates ...")
		return a.renewExternalCustom()
	}

	log.Info("Renewing external certificates ...")
	if !a.cfg.Env.RunInPod {
		log.Debug("Remote execution in progress ...")
		return a.renewExternalNotInPod()
	}
	return a.renewExternalInPod()
}

func (a *renew) renewExternalCustom() error {
	log.Debugf("Read %s", a.cfg.Cert)
	certData, err := utils.ReadFile(a.cfg.Cert)
	if err != nil {
		return err
	}
	log.Debugf("Read %s", a.cfg.Key)
	keyData, err := utils.ReadFile(a.cfg.Key)
	if err != nil {
		return err
	}
	data := make(map[string]string)
	data[a.cfg.ResourceField+".crt"] = string(certData)
	data[a.cfg.ResourceField+".key"] = string(keyData)

	// apply secret
	secrets := strings.Split(a.cfg.Resource, ",")
	for k := range secrets {
		secret := strings.TrimSpace(secrets[k])
		if len(secret) == 0 {
			continue
		}
		log.Debugf("Apply %s in %s", secret, a.cfg.Namespace)
		_, err = a.kube.ApplySecret(a.cfg.Namespace, secret, data)
		if err != nil {
			return errors.Wrapf(err, "apply %s, namespace: %s",secret, a.cfg.Namespace)
		}
	}

	// apply public-ca-certificates configmap
	if len(a.cfg.CACert) > 0 {
		log.Debugf("Read %s", a.cfg.CACert)
		cacertData, err := utils.ReadFile(a.cfg.CACert)
		if err != nil {
			return err
		}
		log.Debugf("Apply %s in %v", ConfigMapNamePublicCA,a.cfg.Namespace)
		newData := make(map[string]string)
		newData["CUS_ca.crt"] = string(cacertData)

		_, err = a.kube.ApplyConfigMap(a.cfg.Namespace, ConfigMapNamePublicCA, newData)
		if err != nil {
			return errors.Wrapf(err, "apply %s, namespace: %s", ConfigMapNamePublicCA, a.cfg.Namespace)
		}
	}

	return nil
}

func (a *renew) renewExternalInPod() error {
	if !strings.Contains(a.cfg.Resource, SecretNameNginxFrontend) {
		log.Debugf("add secret %s ro resource", SecretNameNginxFrontend)
		a.cfg.Resource = fmt.Sprintf("%s,%s", a.cfg.Resource, SecretNameNginxFrontend)
	}
	return a.generator.GenAndDump(&certs.Certificate{
		CN:       a.cfg.Host,
		UintTime: a.cfg.Unit,
		Validity: a.cfg.Validity,
	}, fmt.Sprintf("%s %s", a.cfg.Resource, a.cfg.ResourceField))
}

func (a *renew) renewExternalNotInPod() error {
	return sysc.RenewRERemoteExecution(a.cfg.Env.CDFNamespace, a.cfg.Namespace,
		a.cfg.Unit, a.cfg.Validity, a.cfg.SkipConfirm)
}

func (a *renew) renewInternal() error {
	if a.cfg.Local {
		return a.renewInternalLocal()
	}
	if a.expired {
		err := a.renewInternalExpired()
		if err != nil {
			return err
		}
	}
	return a.renewInternalAvailable()
}

func (a *renew) renewInternalLocal() error {
	return nil
}

func (a *renew) renewInternalExpired() error {
	return nil
}

func (a *renew) renewInternalAvailable() error {
	return nil
}

func getCAs(kube *kube.Client, namespace string) ([][]byte, error) {
	cm, err := kube.GetConfigMap(namespace, ConfigMapNamePublicCA)
	if err != nil {
		return nil, err
	}
	ric, ok := cm.Data[CertNameRIC]
	if !ok {
		return nil, errors.New("RIC ca nil")
	}
	var cas [][]byte
	cas = append(cas, []byte(ric))

	if rid, ok := cm.Data[CertNameRID]; ok {
		log.Debug("got RID ca")
		cas = append(cas, []byte(rid))
	}
	if re, ok := cm.Data[CertNameRE]; ok {
		log.Debug("got RE ca")
		cas = append(cas, []byte(re))
	}
	if cus, ok := cm.Data[CertNameCUS]; ok {
		log.Debug("got CUS ca")
		cas = append(cas, []byte(cus))
	}

	return cas, nil
}

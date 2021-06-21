package action

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

const (
	LabelPrefixDeploymentUuid = "deployments.microfocus.com/deployment-uuid"
	LabelPrefixDeploymentName = "deployments.microfocus.com/deployment-name"
)

type renewSubExternal struct {
	*action

	iValidity IValidity
	isCustom  bool
}

func NewRenewSubExternal(cfg *Configuration, validity IValidity) Interface {
	return &renewSubExternal{
		action:    newAction("renew-sub-external", cfg),
		iValidity: validity,
	}
}

func (a *renewSubExternal) Name() string {
	return a.name
}

func (a *renewSubExternal) Run() error {
	log.Debugf("***** %s Run *****", strings.ToUpper(a.name))
	if !a.cfg.SkipConfirm {
		confirm, err := prompt.Confirm("Are you sure to continue")
		if err != nil {
			return err
		}
		if !confirm {
			return DropError
		}
	}

	var sub Interface
	if len(a.cfg.Cert) > 0 && len(a.cfg.Key) > 0 {
		log.Info("Renewing custom external certificates ...")
		sub = NewRenewSubExternalCustom(a.cfg)
	}

	if !a.cfg.Env.RunInPod {
		log.Info("Renewing external certificates locally ...")
		sub = NewRenewSubExternalNotInPod(a.cfg)
	} else {
		log.Info("Renewing external certificates in Pod ...")
		sub = NewRenewSubExternalInPod(a.cfg)
	}
	return Execute(sub)
}

func (a *renewSubExternal) PreRun() error {
	log.Debugf("***** %s PreRun *****", strings.ToUpper(a.name))

	// ignore checks, if running in a pod
	if a.cfg.Env.RunInPod {
		return nil
	}

	if a.iValidity.server <= 0 {
		log.Error("The internal certificates have already expired.")
		log.Errorf("You should run the '%s/scripts/renewCert --renew' to "+
			"renew the internal certificates firstly.", a.cfg.Env.K8SHome)
		return errors.New("internal certificates expired")
	}

	log.Debug("Checking external RE certificate ...")
	secret, err := a.kube.GetSecret(a.cfg.Namespace, SecretNameNginxDefault)
	if err != nil {
		return err
	}

	if v, ok := secret.Data[DefaultResourceKeyTls+".crt"]; ok {
		recert, err := utils.ParseCrtBytes(v)
		if err != nil {
			return err
		}
		available := utils.CheckCrtValidity(recert)
		if available <= 0 {
			log.Warn("The external certificates have already expired.")
		} else {
			log.Warnf("The external certificates will expire in %d hour(s).", available)
		}
	}

	if len(a.cfg.Cert) > 0 && len(a.cfg.Key) > 0 {
		a.isCustom = true
	} else {
		log.Warn("This command will overwrite the external certificates with self-singed certificates.")
	}

	ns, err := a.kube.GetNamespace(a.cfg.Namespace)
	if err != nil {
		return err
	}
	uuid, ok := ns.Labels[LabelPrefixDeploymentUuid]
	if !ok {
		log.Warn("DEPLOYMENT UUID not found")
		return nil
	}

	log.Debug("Checking if the current deployment is PRIMARY ...")
	log.Debugf("DEPLOYMENT UUID: %s", uuid)
	list, err := a.kube.GetNamespacesWithLabel(fmt.Sprintf("%s=%s", LabelPrefixDeploymentUuid, uuid))
	if err != nil {
		return err
	}

	if len(list.Items) == 0 {
		return errors.Errorf("namespaces with uuid: %s not found", uuid)
	}
	for i := range list.Items {
		if list.Items[i].Name == a.cfg.Env.CDFNamespace {
			a.cfg.Cluster.IsPrimary = true
			break
		}
	}

	return nil
}

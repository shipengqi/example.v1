package sysc

import (
	"fmt"
	"github.com/shipengqi/example.v1/apps/cli/pkg/command"
	"github.com/shipengqi/example.v1/apps/cli/pkg/log"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	DockerSupportVersion = 202108
)

func RestartKubeService(namespace string, version int) error {
	err := RestartNativeService("kubelet")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	err = RestartContainers(namespace, version)
	if err != nil {
		return err
	}
	return nil
}

func RestartNativeService(serviceName string) error {
	log.Infof("Start to restart service: %s.", serviceName)
	err := StopNativeService(serviceName)
	if err != nil {
		return err
	}
	err = StartNativeService(serviceName)
	if err != nil {
		return err
	}
	log.Infof("Restart service %s successfully.", serviceName)
	return nil
}

func StopNativeService(serviceName string) error {
	log.Infof("Stopping service: %s.", serviceName)
	stopCMD := fmt.Sprintf("systemctl stop %s 2>&1", serviceName)
	log.Debugf("exec: %s", stopCMD)
	stdout, _, err := command.Exec(stopCMD)
	if err != nil {
		return errors.Wrap(err, stdout)
	}
	return nil
}

func StartNativeService(serviceName string) error {
	log.Infof("Starting service: %s.", serviceName)
	stdout, _, err := command.Exec("systemctl daemon-reload")
	if err != nil {
		return errors.Wrap(err, stdout)
	}
	startCMD := fmt.Sprintf("systemctl start %s 2>&1", serviceName)
	log.Debugf("exec: %s", startCMD)
	stdout, _, err = command.Exec(startCMD)
	if err != nil {
		return errors.Wrap(err, stdout)
	}
	return nil
}

func RestartContainers(namespace string, version int) error {
	var restartCMD string
	restartCMD = fmt.Sprintf("docker restart $(docker ps | grep %s | " +
		"sed '/\\/pause/d' | awk '{print $1}')", namespace)
	if version > DockerSupportVersion {
		// TODO if ps command does not fetch any resources, skip
		restartCMD = fmt.Sprintf("crictl stop $(crictl ps --label io.kubernetes.pod.namespace=%s | " +
			"grep '^[^CONTAINER]' | awk '{print $1}')", namespace)
	}

	log.Debugf("exec: %s", restartCMD)
	log.Infof("Start to restart all containers in namespace: %s.", namespace)
	err := command.ExecSync(restartCMD)
	if err != nil {
		return errors.Wrap(err, "restart containers")
	}
	log.Info("Restart all containers successfully.")
	return nil
}

func ParseVaultToken(encryptedToken, passphrase, tokenEncKey, tokenEncIv string) (string, error) {
	cmd := fmt.Sprintf(`echo %s ` +
		`| openssl aes-256-cbc -md sha256 -a -d ` +
		`-pass pass:"%s" -K %s -iv %s`,
		encryptedToken,
		passphrase,
		tokenEncKey,
		tokenEncIv)
	log.Debugf("exec: %s", cmd)
	stdout, _, err := command.Exec(cmd)
	if err != nil {
		return "", errors.Wrap(err, stdout)
	}
	token := strings.Trim(stdout, "\n")
	return token, nil
}

func RenewRERemoteExecution(cdfNamespace, namespace, unit, resource, field string, primary bool, V int) error {
	// Todo filter the running pod
	template := `kubectl exec -n %s $(kubectl get po -n %s ` +
		`-l 'deployments.microfocus.com/component=itom-vault' ` +
		`-o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}') ` +
		`-- /renewCert --renew -t external -V %d --unit-time %s -n %s --cdf-namespace %s ` +
		`--resource %s --field %s -y`
	primaryNs := namespace
	if primary {
		primaryNs = cdfNamespace
		template += " --primary"
	}
	cmd := fmt.Sprintf(template,
		primaryNs,
		primaryNs,
		V,
		unit,
		namespace,
		cdfNamespace,
		resource,
		field)
	log.Debugf("exec: %s", cmd)
	err := command.ExecSync(cmd)
	if err != nil {
		return err
	}

	return nil
}

func Hostname() (string, error) {
	cmd := `systemctl show --property=ExecStart kubelet | grep "hostname-override"`
	log.Debugf("exec: %s", cmd)
	stdout, _, err := command.Exec(cmd)
	if err != nil {
		return "", errors.Wrap(err, stdout)
	}

	r, err := regexp.Compile(`hostname-override=.*\s-`)
	if err != nil {
		return "", err
	}

	override := r.FindString(stdout)
	words := strings.Split(override, "=")
	words2 := strings.Split(words[1], " ")

	return strings.ToLower(words2[0]), nil
}

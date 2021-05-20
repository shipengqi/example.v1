package sysc

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/command"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func RestartKubeService(namespace string) error {
	err := RestartNativeService("kubelet")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	err = RestartContainersWithNamespace(namespace)
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

func RestartContainersWithNamespace(namespace string) error {
	restartCMD := fmt.Sprintf("docker restart $(docker ps | grep %s | sed '/\\/pause/d' | awk '{print $1}')", namespace)
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

func RenewREInContainer(cdfNamespace, namespace, unit string, V int, confirm bool) error {
	cmd := fmt.Sprintf(`kubectl exec -n %s $(kubectl get po -n %s ` +
		`-l 'deployments.microfocus.com/component=itom-vault' ` +
		`-o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}') ` +
		`-- /renewCert --renew -t external -V %d --unit-time %s -n %s --cdf-namespace %s -y %v`,
		namespace,
		namespace,
		V,
		unit,
		namespace,
		cdfNamespace,
		confirm)
	log.Debugf("exec: %s", cmd)
	err := command.ExecSync(cmd)
	if err != nil {
		return err
	}

	return nil
}

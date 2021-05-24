/*
* © Copyright 2017-2019 Micro Focus or one of its affiliates.
*
* The only warranties for products and services of Micro Focus and its
* affiliates and licensors ("Micro Focus") are as may be set forth in
* the express warranty statements accompanying such products and services.
* Nothing herein should be construed as constituting an additional warranty.
* Micro Focus shall not be liable for technical or editorial errors or
* omissions contained herein. The information contained herein is subject
* to change without notice.
*
* Except as specifically indicated otherwise, this document contains
* confidential information and a valid license is required for possession,
* use or copying. If this work is provided to the U.S. Government, consistent
* with FAR 12.211 and 12.212, Commercial Computer Software, Computer Software
* Documentation, and Technical Data for Commercial Items are licensed to the
* U.S. Government under vendor's standard commercial license.
 */

package env

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ITOMCDFEnvFile        = "/etc/profile.d/itom-cdf.sh"
	EnvKeyK8SHome         = "K8S_HOME"
	EnvKeyCDFNamespace    = "CDF_NAMESPACE"
	EnvKeyRuntimeDataHome = "RUNTIME_CDFDATA_HOME"
	EnvKeyVaultRole       = "CERTIFICATE_ROLE"
)

const (
	_defaultVaultRole       = "coretech"
	_defaultCDFNamespace    = "core"
	_defaultK8SHome         = "/opt/kubernetes"
	_defaultRunTimeDataPath = "/opt/kubernetes/data"
	_defaultK8STokenFile    = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Settings struct {
	K8SHome            string
	CDFNamespace       string
	RuntimeCDFDataHome string
	SSLPath            string
	VaultRole          string
	RunOnMaster        bool
	RunInPod           bool
}

func New() (*Settings, error) {
	envs := &Settings{
		K8SHome:            envOr(EnvKeyK8SHome, _defaultK8SHome),
		CDFNamespace:       envOr(EnvKeyCDFNamespace, _defaultCDFNamespace),
		RuntimeCDFDataHome: envOr(EnvKeyRuntimeDataHome, _defaultRunTimeDataPath),
		VaultRole:          envOr(EnvKeyVaultRole, _defaultVaultRole),
		SSLPath:            "",
		RunOnMaster:        false,
		RunInPod:           false,
	}

	envs.RunInPod = inPod()
	if envs.RunInPod {
		return envs, nil
	}

	values, err := retrieveEnv(ITOMCDFEnvFile, EnvKeyK8SHome, EnvKeyCDFNamespace)
	if err != nil {
		return nil, err
	}
	if v, ok := values[EnvKeyCDFNamespace]; ok {
		log.Debugf("got env: %s, value: %s ", EnvKeyCDFNamespace, v)
		envs.CDFNamespace = v
	}
	if v, ok := values[EnvKeyK8SHome]; ok {
		log.Debugf("got env: %s, value: %s ", EnvKeyK8SHome, v)
		envs.K8SHome = v

		values, err = retrieveEnv(fmt.Sprintf("%s/bin/env.sh", v), EnvKeyRuntimeDataHome)
		if err != nil {
			return envs, err
		}
		if dataHome, ok := values[EnvKeyRuntimeDataHome]; ok {
			log.Debugf("got env: %s, value: %s ", EnvKeyRuntimeDataHome, dataHome)
			envs.RuntimeCDFDataHome = dataHome
		}
	}

	envs.RunOnMaster = onMasterNode(envs.RuntimeCDFDataHome)
	envs.SSLPath = filepath.Join(envs.K8SHome, "ssl")

	return envs, nil
}

// ----------------------------------------------------------------------------
// Helpers...

func retrieveEnv(filePath string, keys ...string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	mappings := make(map[string]string)
	for k := range keys {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lineText := scanner.Text()
			if strings.Contains(lineText, fmt.Sprintf("%s=", keys[k])) {
				ss := strings.Split(lineText, "=")
				mappings[keys[k]] = ss[len(ss)-1]
				break
			}
		}
	}

	return mappings, nil
}

func isExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func onMasterNode(dataHome string) bool {
	etcdDir := dataHome + "/etcd"
	if isExist(etcdDir) {
		files, _ := ioutil.ReadDir(etcdDir)
		if len(files) == 0 {
			return false
		}
		return true
	}
	return false
}

func inPod() bool {
	return isExist(_defaultK8STokenFile)
}

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

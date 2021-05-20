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

package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ITOMCDFEnvFile = "/etc/profile.d/itom-cdf.sh"

	EnvKeyK8SHome         = "K8S_HOME="
	EnvKeyCDFNamespace    = "CDF_NAMESPACE="
	EnvKeyRuntimeDataHome = "RUNTIME_CDFDATA_HOME="
)

const (
	_defaultCDFNamespace    = "core"
	_defaultK8SHome         = "/opt/kubernetes"
	_defaultRunTimeDataPath = "/opt/kubernetes/data"
	_defaultK8STokenFile    = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type Envs struct {
	K8SHome            string
	CDFNamespace       string
	RuntimeCDFDataHome string
	RunOnMaster        bool
	RunInPod           bool
}

func InitEnvs() (*Envs, error) {
	envs := &Envs{
		K8SHome:            _defaultK8SHome,
		CDFNamespace:       _defaultCDFNamespace,
		RuntimeCDFDataHome: _defaultRunTimeDataPath,
		RunOnMaster:        false,
		RunInPod:           false,
	}

	values, err := retrieveEnv(ITOMCDFEnvFile, EnvKeyK8SHome, EnvKeyCDFNamespace)
	if err != nil {
		return nil, errors.Wrapf(err, "open %s", ITOMCDFEnvFile)
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
			return envs, errors.Wrap(err, "open env.sh")
		}
		if dataHome, ok := values[EnvKeyRuntimeDataHome]; ok {
			log.Debugf("got env: %s, value: %s ", EnvKeyRuntimeDataHome, dataHome)
			envs.RuntimeCDFDataHome = dataHome
		}
	}

	envs.RunOnMaster = onMasterNode(envs.RuntimeCDFDataHome)
	envs.RunInPod = inPod()

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
			if strings.Contains(lineText, keys[k]) {
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

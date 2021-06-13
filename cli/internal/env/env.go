package env

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shipengqi/example.v1/cli/pkg/command"
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
	_defaultCDFVersion      = 202108
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
	Version            int
}

func New() (*Settings, error) {
	envs := &Settings{
		K8SHome:            envOr(EnvKeyK8SHome, _defaultK8SHome),
		CDFNamespace:       envOr(EnvKeyCDFNamespace, _defaultCDFNamespace),
		RuntimeCDFDataHome: envOr(EnvKeyRuntimeDataHome, _defaultRunTimeDataPath),
		VaultRole:          envOr(EnvKeyVaultRole, _defaultVaultRole),
		SSLPath:            "",
		RunOnMaster:        true,
		RunInPod:           false,
		Version:            _defaultCDFVersion,
	}

	envs.RunInPod = inPod()
	if envs.RunInPod {
		return envs, nil
	}

	values, err := retrieveEnv(ITOMCDFEnvFile, EnvKeyK8SHome, EnvKeyCDFNamespace)
	if err != nil {
		log.Warnf("retrieveEnv(): %v", err)
	}
	if v, ok := values[EnvKeyCDFNamespace]; ok {
		log.Debugf("got env: %s, value: %s ", EnvKeyCDFNamespace, v)
		envs.CDFNamespace = v
	}
	if v, ok := values[EnvKeyK8SHome]; ok {
		log.Debugf("got env: %s, value: %s ", EnvKeyK8SHome, v)
		envs.K8SHome = v

		envfile := fmt.Sprintf("%s/bin/env.sh", v)
		values, err = retrieveEnv(envfile, EnvKeyRuntimeDataHome)
		if err != nil {
			log.Warnf("retrieveEnv(): %v", err)
		}
		if dataHome, ok := values[EnvKeyRuntimeDataHome]; ok {
			log.Debugf("got env: %s, value: %s ", EnvKeyRuntimeDataHome, dataHome)
			envs.RuntimeCDFDataHome = dataHome
		}
	}

	envs.RunOnMaster = onMasterNode(envs.RuntimeCDFDataHome)
	envs.SSLPath = filepath.Join(envs.K8SHome, "ssl")

	version, _, err := command.Exec(fmt.Sprintf("cat %s/version.txt | awk -F . '{print $1$2}'", envs.K8SHome))
	if err != nil {
		log.Warnf("open version.txt: %v", err)
		return envs, err
	}
	vi, err := strconv.Atoi(strings.TrimSpace(version))
	if err != nil {
		log.Warnf("strconv.Atoi(): %v", err)
	} else {
		envs.Version = vi
	}

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

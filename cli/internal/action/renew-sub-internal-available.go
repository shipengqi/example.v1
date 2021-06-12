package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/certs/infra"
	"github.com/shipengqi/example.v1/cli/internal/node"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/prompt"
)

var (
	maxRetryTimes     = 3
	currentRetryTimes = 0
)

type credential struct {
	mode    string
	user    string
	passwd  string
	keyfile string
	singer  ssh.Signer
}

type renewSubInternalAvailable struct {
	*action

	generator certs.Generator
}

func NewRenewSubInternalAvailable(cfg *Configuration) Interface {
	c := &renewSubInternalAvailable{
		action: newAction("renew-sub-internal-available", cfg),
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
	log.Debug("gen cert for all nodes")

	for _, v := range nodes {
		if v.Address == a.cfg.Cluster.FirstMasterNode {
			v.First = true
		}
		// avoid duplicate renew cert on master node
		// if v.Address == currentHost || v.Address == currentIP {}
		err = a.iterate(v.Address, v.Master, false, a.generator)
		if err != nil {
			return err
		}
	}

	log.Info("Generate certificates successfully.")

	return a.distribute(nodes)
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

func (a *renewSubInternalAvailable) overwriteAllCerts(n node.Node) (int, error) {
	for _, v := range CertificateSet {
		if !v.CanDep(n.Master) {
			continue
		}

		srcCrt := path.Join(a.cfg.OutputDir, fmt.Sprintf("%s-%s.crt", n.Address, v.Name))
		srcKey := path.Join(a.cfg.OutputDir, fmt.Sprintf("%s-%s.key", n.Address, v.Name))
		dstCrt := path.Join(a.cfg.Env.SSLPath, fmt.Sprintf("%s.crt", v.Name))
		dstKey := path.Join(a.cfg.Env.SSLPath, fmt.Sprintf("%s.key", v.Name))
		status, err := n.Scp(srcCrt, dstCrt)
		if err != nil {
			return status, err
		}
		log.Debugf("distribute %s to %s.", dstCrt, n.Address)
		log.Debugf("remove %s", srcCrt)
		_ = os.Remove(srcCrt)

		status, err = n.Scp(srcKey, dstKey)
		if err != nil {
			return status, err
		}
		log.Debugf("distribute %s to %s.", dstKey, n.Address)
		log.Debugf("remove %s", srcKey)
		_ = os.Remove(srcKey)

	}
	log.Infof("Distribute certificates to %s successfully.", n.Address)
	return 0, nil
}

func (a *renewSubInternalAvailable) distribute(nodes []node.Node) error {
	if !a.cfg.SkipConfirm {
		confirm, err := prompt.Confirm("Do you want to distribute certificates to all the nodes")
		if err != nil {
			return err
		}
		if !confirm {
			log.Infof("You can distribute the certificates under %s manually.", a.cfg.OutputDir)
			log.Infof("After that, run '%s/scripts/renewCert --apply' "+
				"one each node one by one so that the certificates can take effect.", a.cfg.Env.K8SHome)
			return nil
		}
	}

	c, err := collectCredential(a.cfg.Options, nodes[0])
	if err != nil {
		return err
	}
	log.Debug("***** distributing *****")
	// This part is try to connect all the nodes.
	// nodes left to connect
	lastNodes := nodes

	// the list of successfully connected nodes and those are not
	var connectedNodes, unconnectedNodes []node.Node
	unconnectedNodes = nil
	log.Info("Connecting ...")
	// try to connect each node
	for _, v := range lastNodes {
		log.Debugf("Try to connect to %s", v.Address)
		status, err := v.Connect(c.user, c.passwd, c.singer)
		v.Err = err
		if status == 1 || status == 2 {
			unconnectedNodes = append(unconnectedNodes, v)
		} else {
			connectedNodes = append(connectedNodes, v)
		}
	}

	// show nodes connected successfully
	if len(connectedNodes) > 0 {
		log.Info("[Successful connection nodes]:")
		for _, v := range connectedNodes {
			log.Infof("    - %s", v.Address)
		}
	}

	// show nodes failed to connect
	if len(unconnectedNodes) > 0 {
		log.Warn("[Failed connection nodes]:")
		for _, v := range unconnectedNodes {
			log.Warnf("    - %s", v.Address)
			log.Debugf("    Err: %s", v.Err)
		}
		log.Warnf("Failed to connect all of the nodes. "+
			"Please distribute the certificates under %s manually.", a.cfg.OutputDir)
		log.Warnf("And then run '%s/scripts/renewCert --apply' "+
			"one each node one by one so that the certificates can take effect.", a.cfg.Env.K8SHome)
		return errors.New("unconnected")
	}

	// if no success connection, exit
	if len(connectedNodes) < 1 {
		log.Warn("None of the nodes can be connected.")
		return errors.New("unconnected")
	}

	// This part is to distribute certificates on all of the successfully connect nodes.
	log.Info("Start to distribute certificates.")
	var unknownErrNodes, successfulNodes, authErrNodes []node.Node
	for _, v := range connectedNodes {
		status, err := a.overwriteAllCerts(v)
		v.Err = err
		if status == 1 {
			authErrNodes = append(authErrNodes, v)
		} else if status == 2 {
			unknownErrNodes = append(unknownErrNodes, v)
		} else {
			successfulNodes = append(successfulNodes, v)
		}
	}
	if len(authErrNodes) > 0 {
		log.Warn("[Authentication error nodes]:")
		for _, v := range authErrNodes {
			log.Warnf("    - %s", v.Address)
			log.Debugf("    Err: %s", v.Err)
		}
	}

	if len(unknownErrNodes) > 0 {
		log.Warn("[Unknown failed nodes]:")
		for _, v := range unknownErrNodes {
			log.Warnf("    - %s", v.Address)
			log.Debugf("    Err: %s", v.Err)
		}
	}

	if len(successfulNodes) > 0 {
		log.Info("[Successful distribution nodes]:")
		for _, v := range successfulNodes {
			log.Infof("    - %s", v.Address)
		}
	}

	return a.applyOnNodes(successfulNodes)
}

func (a *renewSubInternalAvailable) applyOnNodes(nodes []node.Node) error {
	if !a.cfg.SkipConfirm {
		confirm, err := prompt.Confirm("Do you want to apply certificates for successful nodes")
		if err != nil {
			return err
		}
		if !confirm {
			log.Infof("Run '%s/bin/renewCert --apply' one each node one by one to "+
				"make the certificates take effect.", a.cfg.Env.K8SHome)
			return nil
		}
	}

	log.Debug("***** applying *****")
	if len(nodes) < 1 {
		log.Debug("Nodes count: 0")
		return nil
	}

	for _, v := range nodes {
		err := v.Exec(path.Join(a.cfg.Env.K8SHome, "scripts/renewCert --apply --remote"))
		if err != nil {
			log.Errorf("Apply certificates on %s, ERR: %v", v.Address, err)
			log.Error("")
		} else {
			log.Infof("Apply certificates on %s successfully.", v.Address)
			log.Info("")
		}
	}

	return nil
}

func collectCredential(flags *Options, n node.Node) (*credential, error) {
	if currentRetryTimes >= maxRetryTimes {
		return nil, errors.New("authenticate failed")
	}
	c := &credential{
		mode:    "",
		user:    "",
		passwd:  "",
		keyfile: "",
	}

	if flags.SSHKey != "" || (flags.Password != "" && flags.Username != "") {
		c.user = flags.Username
		c.passwd = flags.Password
		c.keyfile = flags.SSHKey
		return c, nil
	}

	mode, err := prompt.Select(
		fmt.Sprintf("SSH mode for %s and the rest nodes", n.Address),
		[]string{"password", "private key"},
	)
	if err != nil {
		return nil, err
	}

	user, err := prompt.Input(fmt.Sprintf("Please input node 'user' for node: %s \n", n.Address))
	if err != nil {
		return nil, err
	}
	c.mode = mode
	c.user = user
	if c.mode == "password" {
		pass, err := prompt.Password(fmt.Sprintf("Please input node 'password' for node: %s \n", n.Address))
		if err != nil {
			return nil, err
		}
		c.passwd = pass
	} else if c.mode == "private key" {
		file, err := prompt.Input(fmt.Sprintf("Please input node private key path for node %s \n", n.Address))
		if err != nil {
			return nil, err
		}
		c.keyfile = file
	}

	if c.keyfile != "" {
		s, err := getSigner(c.keyfile)
		if err != nil {
			return nil, errors.Wrap(err, "get signer")
		}
		c.singer = s
	}
	log.Info("Testing ...")
	status, err := n.Connect(c.user, c.passwd, c.singer)
	if status == 1 {
		log.Warnf("Authenticate failed - %s, please try again.", n.Address)
		currentRetryTimes++
		return collectCredential(flags, n)
	}
	return c, nil
}

func getSigner(path string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "read private key")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "parse private key")
	}
	return signer, nil

}

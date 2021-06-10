package node

import (
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/shipengqi/example.v1/cli/pkg/log"
	ssh2 "github.com/shipengqi/example.v1/cli/pkg/ssh"
)

type Node struct {
	Address string
	User    string
	Pass    string
	Master  bool
	First   bool
	Err     error
	Signer  ssh.Signer
	client  ssh2.Client
}

func (n *Node) Connect(user, pass string, signer ssh.Signer) (int, error) {
	if pass != "" {
		n.client = ssh2.New(n.Address+":22", user, pass, nil)
	} else {
		n.client = ssh2.New(n.Address+":22", user, "", signer)
	}

	err := n.client.Connect()
	if err != nil {
		if strings.Contains(err.Error(), "authenticate") {
			return 1, err
		} else {
			return 2, err
		}
	}
	n.User = user
	n.Pass = pass
	n.Signer = signer
	return 0, nil
}

func (n *Node) Exec(command string) error {
	log.Debugf("exec %s", command)
	_, err := n.Connect(n.User, n.Pass, n.Signer)
	if err != nil {
		return errors.Wrap(err, "connect "+n.Address)
	}
	err = n.client.OutputSync(command)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Scp(src, dst string) (int, error) {
	log.Debugf("scp %s to %s", src, dst)
	status, err := n.Connect(n.User, n.Pass, n.Signer)
	if err != nil {
		return status, errors.Wrap(err, "connect "+n.Address)
	}
	err = n.client.Scp(src, dst)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

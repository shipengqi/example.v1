package ssh

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/shipengqi/example.v1/cli/pkg/log"
)

type Client struct {
	addr   string
	user   string
	auths  []ssh.AuthMethod
	client *ssh.Client
}

func New(addr, user, password string, signer ssh.Signer) Client {
	var authMethod []ssh.AuthMethod
	if password != "" {
		authMethod = append(authMethod, ssh.Password(password))
	}
	if signer != nil {
		authMethod = append(authMethod, ssh.PublicKeys(signer))
	}
	sshClient := Client{addr: addr, auths: authMethod, user: user}
	return sshClient
}

func (s *Client) Connect() error {
	config := &ssh.ClientConfig{
		User: s.user,
		Auth: s.auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * 10,
	}

	client, err := ssh.Dial("tcp", s.addr, config)
	if err != nil {
		return errors.Wrapf(err, "connect %s", s.addr)
	}
	s.client = client
	return nil
}

func (s *Client) Output(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", errors.Wrap(err, "create ssh session")
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", errors.Wrap(err, "run cmd")
	}
	return string(output), nil
}

func (s *Client) OutputSync(cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	err = session.Start(cmd)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, _, err2 := reader.ReadLine()
		if err2 != nil || io.EOF == err2 {
			break
		}
		log.Info(string(line))
	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *Client) Scp(sourceFile string, etargetFile string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return errors.Wrap(err, "create ssh session")
	}

	defer session.Close()

	targetFile := filepath.Base(etargetFile)

	src, srcErr := os.Open(sourceFile)
	if srcErr != nil {
		return errors.Wrap(err, "open source file")
	}

	srcStat, statErr := src.Stat()

	if statErr != nil {
		return statErr
	}

	go func() {
		w, _ := session.StdinPipe()

		fmt.Fprintln(w, "C0644", srcStat.Size(), targetFile)

		if srcStat.Size() > 0 {
			io.Copy(w, src)
			fmt.Fprint(w, "\x00")
			w.Close()
		} else {
			fmt.Fprint(w, "\x00")
			w.Close()
		}
	}()

	if err := session.Run(fmt.Sprintf("scp -tr %s", etargetFile)); err != nil {
		return errors.Wrap(err, "run scp cmd")
	}
	return nil
}

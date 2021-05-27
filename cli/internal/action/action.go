package action

import (
	"net"
	"strings"

	"github.com/shipengqi/example.v1/cli/pkg/log"
)


type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Execute() error
}

type action struct {
	name string
	cfg  *Configuration
}

func (a *action) Name() string {
	return "[action]"
}

func (a *action) PreRun() error {
	log.Debugf("====================    %s PreRun    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) Run() error {
	log.Debugf("====================   %s Run    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) PostRun() error {
	log.Debugf("====================   %s PostRun    ====================", strings.ToUpper(a.name))
	return nil
}

func (a *action) Execute() error {
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

// ----------------------------------------------------------------------------
// Helpers...

func parseSan(san string) ([]string, []net.IP, string) {
	if len(san) == 0 {
		return nil, nil, ""
	}

	var svcIp string
	dns := make([]string, 0)
	ips := make([]net.IP, 0)

	subs := strings.Split(san, ",")
	for _, sub := range subs {
		if strings.HasPrefix(sub, "DNS:") {
			dns = append(dns, sub[4:])
		}
		if strings.HasPrefix(sub, "IP:") {
			ips = append(ips, net.ParseIP(sub[3:]))
		}
		if strings.HasPrefix(sub, "K8SSVCIP:") {
			svcIp = sub[9:]
		}
	}

	return dns, ips, svcIp
}

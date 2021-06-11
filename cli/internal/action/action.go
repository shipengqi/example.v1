package action

import (
	"crypto"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

const (
	ConfigMapNameCDFCluster = "cdf-cluster-host"
	ConfigMapNameCDF        = "cdf"
	ConfigMapNamePublicCA   = "public-ca-certificates"
	SecretNameNginxDefault  = "nginx-default-secret"
	SecretNameNginxFrontend = "nginx-frontend-secret"
	SecretNameK8SRootCert   = "k8s-root-cert"
	ResourceKeyRECert       = "RE_ca.crt"
	ResourceKeyRICCert      = "RIC_ca.crt"
	ResourceKeyRIDCert      = "RID_ca.crt"
	ResourceKeyCUSCert      = "CUS_ca.crt"
	ResourceKeyCAKey        = "ca.key"
)

const (
	DefaultResourceKeyTls = "tls"
)

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
}

type action struct {
	name string
	cfg  *Configuration
	kube *kube.Client
}

func newAction(name string, cfg *Configuration) *action {
	c := &action{
		name: name,
		cfg:  cfg,
	}

	kclient, err := kube.New(cfg.Kube)
	if err != nil {
		panic(err)
	}

	c.kube = kclient

	return c
}

func (a *action) Name() string {
	return "[action]"
}

func (a *action) PreRun() error {
	log.Debugf("***** [%s] PreRun *****", strings.ToUpper(a.name))
	a.cfg.Debug()
	return nil
}

func (a *action) Run() error {
	log.Debugf("***** [%s] Run *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) PostRun() error {
	log.Debugf("***** [%s] PostRun *****", strings.ToUpper(a.name))
	return nil
}

func Execute(a Interface) error {
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

func (a *action) iterate(address string, master, overwrite bool, generator certs.Generator) error {
	dnsnames, ipaddrs, cn := a.combineSubject(address, master)
	log.Debugf("initial cert DNS: %s", dnsnames)
	log.Debugf("initial cert IPs: %s", ipaddrs)
	log.Debugf("initial cert CN: %s", cn)

	for _, v := range CertificateSet {
		if !v.CanDep(master) {
			continue
		}

		log.Debug("********** START **********")
		v.Host = address
		v.Validity = a.cfg.Validity
		v.UintTime = a.cfg.Unit
		v.Overwrite = overwrite

		dns := make([]string, len(dnsnames))
		ips := make([]net.IP, len(ipaddrs))
		copy(dns, dnsnames)
		copy(ips, ipaddrs)

		v.CombineServerSan(dns, ips, cn, a.cfg.ServerCertSan, a.cfg.Cluster.KubeServiceIP)

		log.Debugf("cert DNS: %s", v.DNSNames)
		log.Debugf("cert IPs: %s", v.IPs)
		log.Debugf("cert CN: %s", v.CN)

		err := generator.GenAndDump(v.Certificate, a.cfg.OutputDir)
		if err != nil {
			return err
		}

		log.Debug("********** END **********")
		log.Debug("")
	}

	return nil
}

func (a *action) iterateSecrets(generator certs.Generator) error {
	for _, v := range CertificateSecretSet {
		log.Debugf("gen secret: %s, cert: %s", v.Secret, v.Name)
		v.Validity = a.cfg.Validity
		v.UintTime = a.cfg.Unit
		log.Debugf("cert validity: %d, unit: %s", v.Validity, v.UintTime)

		secretName, secretNs := parseSecretName(v.Secret)
		if len(secretNs) == 0 {
			secretNs = a.cfg.Env.CDFNamespace
		}
		if v.IsKubeRegistryCert() {
			cn := fmt.Sprintf("%s.%s", v.Name, secretNs)
			v.CN = cn
			v.DNSNames = []string{
				"localhost",
				cn,
			}
		}
		log.Debugf("cert DNS: %s", v.DNSNames)
		log.Debugf("cert IPs: %s", v.IPs)
		log.Debugf("cert CN: %s", v.CN)

		crt, key, err := generator.Gen(v.Certificate)
		if err != nil {
			return err
		}

		data := make(map[string][]byte)
		data[v.Name+".crt"] = crt
		data[v.Name+".key"] = key
		_, err = a.kube.ApplySecretBytes(secretNs, secretName, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *action) combineSubject(address string, master bool) ([]string, []net.IP, string) {
	cn := address
	var dnsNames []string
	var ips []net.IP

	if !master {
		if utils.IsIPV4(address) {
			cn = fmt.Sprintf("host-%s", address)
		}

		return dnsNames, ips, cn
	}

	if a.cfg.Cluster.VirtualIP != "" {
		ips = append(ips, net.ParseIP(a.cfg.Cluster.VirtualIP))
	}
	if a.cfg.Cluster.LoadBalanceIP != "" {
		lbIp := net.ParseIP(a.cfg.Cluster.LoadBalanceIP)
		if lbIp != nil {
			ips = append(ips, lbIp)
		} else {
			dnsNames = append(dnsNames, a.cfg.Cluster.LoadBalanceIP)
		}
	}
	if utils.IsIPV4(address) {
		cn = fmt.Sprintf("host-%s", address)
		ips = append(ips, net.ParseIP(address))
	} else {
		dnsNames = append(dnsNames, address)
		shortNames := strings.Split(address, ".")
		// Add first short name
		if len(shortNames) > 0 {
			dnsNames = append(dnsNames, shortNames[0])
		}

		// Add master node ip
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP)
						break
					}
				}
			}
		}
	}

	return dnsNames, ips, cn
}

func (a *action) parseCAKey() (crypto.PrivateKey, error) {
	if len(a.cfg.CAKey) > 0 && utils.IsExist(a.cfg.CAKey) {
		log.Debugf("ParseKey(): %s", a.cfg.CAKey)
		return utils.ParseKey(a.cfg.CAKey)
	} else {
		secret, err := a.kube.GetSecret(NamespaceKubeSystem, SecretNameK8SRootCert)
		if err != nil {
			return nil, err
		}
		if v, ok := secret.Data[ResourceKeyCAKey]; ok {
			return utils.ParseKeyBytes(v, false)
		}
	}
	return nil, errors.New("ca key is nil")
}

// ----------------------------------------------------------------------------
// Helpers...

func parseSan(san string) ([]string, []net.IP, net.IP) {
	if len(san) == 0 {
		return nil, nil, nil
	}

	var svcIp net.IP
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
			svcIp = net.ParseIP(sub[9:])
		}
	}

	return dns, ips, svcIp
}

func parseSecretName(resource string) (name, namespace string) {
	if len(resource) == 0 {
		return
	}

	subs := strings.Split(resource, ".")
	if len(subs) == 1 {
		return subs[0], ""
	}
	if len(subs) >= 2 {
		return subs[0], subs[1]
	}
	return
}

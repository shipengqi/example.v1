package certs

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	rd "math/rand"
	"net"
	"time"

	"github.com/shipengqi/example.v1/cli/internal/types"
)

const BaseDuration = 24 * 60 * 60 * 1000 * 1000 * 1000

type Certificate struct {
	Name          string
	CN            string
	Host          string // node host, used when distributing certificates
	UintTime      string // used for testing
	Overwrite     bool   // local or expired certificates need to overwrite
	IsCA          bool
	Validity      int
	KeyUsage      x509.KeyUsage
	Organizations []string
	DNSNames      []string
	IPs           []net.IP
	ExtKeyUsages  []x509.ExtKeyUsage
}

func (c *Certificate) Gen() *x509.Certificate {
	subject := pkix.Name{
		CommonName: c.CN,
	}
	subject.Organization = c.Organizations
	duration := time.Duration(c.validity())
	obj := &x509.Certificate{
		SerialNumber:          big.NewInt(rd.Int63()),
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(duration),
		BasicConstraintsValid: true,
		IsCA:                  c.IsCA,
		KeyUsage:              c.KeyUsage,
		ExtKeyUsage:           c.ExtKeyUsages,
	}

	if c.DNSNames != nil {
		obj.DNSNames = c.DNSNames
	}

	if c.IPs != nil {
		obj.IPAddresses = c.IPs
	}

	return obj
}

func (c *Certificate) validity() int {
	d := BaseDuration
	if c.UintTime == types.CertUnitTimeMinute {
		d = 60 * 1000 * 1000 * 1000
	}
	return c.Validity * d
}

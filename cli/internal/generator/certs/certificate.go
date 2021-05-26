package certs

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	rd "math/rand"
	"net"
	"time"
)

type Certificate struct {
	CN            string
	UintTime      string
	IsCA          bool
	Period        int
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
	duration := time.Duration(c.Period)
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

package certs

import (
	"crypto/x509"
	"net"
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

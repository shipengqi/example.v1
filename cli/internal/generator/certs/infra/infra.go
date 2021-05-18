package infra

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	rd "math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/keys"
)

func New() certs.Generator {
	return &generator{}
}

type generator struct {
	rootCA  *x509.Certificate
	rootKey crypto.PrivateKey
	keys    keys.Generator
}

func (g *generator) Gen(c *certs.Certificate) (cert, key string, err error) {
	privk, err := g.keys.Gen()
	if err != nil {
		return "", "", err
	}
	pubk := privk.Public()

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

	objBytes, err := x509.CreateCertificate(rand.Reader, obj, g.rootCA, pubk, g.rootKey)
	if err != nil {
		return "", "", errors.Wrapf(err, "x509.CreateCertificate")
	}

	objPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: objBytes,
	}

	return string(pem.EncodeToMemory(objPem)), string(g.keys.Encode(privk)), nil
}

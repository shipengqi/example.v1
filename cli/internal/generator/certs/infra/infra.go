package infra

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"github.com/pkg/errors"
	"io/ioutil"

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

func (g *generator) Gen(c *certs.Certificate) (cert, key []byte, err error) {
	privk, err := g.keys.Gen()
	if err != nil {
		return nil, nil, err
	}
	pubk := privk.Public()

	objBytes, err := x509.CreateCertificate(rand.Reader, c.Gen(), g.rootCA, pubk, g.rootKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "x509.CreateCertificate")
	}

	objPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: objBytes,
	}

	return pem.EncodeToMemory(objPem), g.keys.Encode(privk), nil
}

func (g *generator) Dump(certName, keyName, secret string, cert, key []byte) (err error) {
	err = ioutil.WriteFile(certName, cert, 0400)
	if err != nil {
		return errors.Wrapf(err, "write %s", certName)
	}
	err = ioutil.WriteFile(keyName, key, 0400)
	if err != nil {
		return errors.Wrapf(err, "write %s", keyName)
	}

	return
}

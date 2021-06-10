package infra

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/generator/keys"
	"github.com/shipengqi/example.v1/cli/internal/generator/keys/rsa"
	"github.com/shipengqi/example.v1/cli/internal/utils"
	"github.com/shipengqi/example.v1/cli/pkg/log"
)

func New(cacrt string, cakey crypto.PrivateKey) (certs.Generator, error) {
	var ca  *x509.Certificate
	var err error


	ca, err = utils.ParseCrt(cacrt)
	if err != nil {
		return nil, err
	}
	if cakey == nil {
		return nil, errors.New("ca key is nil")
	}

	return &generator{
		rootCA:  ca,
		rootKey: cakey,
		keys:    rsa.New(),
	}, nil
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

func (g *generator) GenAndDump(c *certs.Certificate, output string) (err error) {
	crtName := path.Join(output, fmt.Sprintf("%s-%s.crt", c.CN, c.Name))
	keyName := path.Join(output, fmt.Sprintf("%s-%s.key", c.CN, c.Name))

	if c.Overwrite {
		crtName = path.Join(output, fmt.Sprintf("%s.crt", c.Name))
		keyName = path.Join(output, fmt.Sprintf("%s.key", c.Name))
		if utils.IsExist(crtName) {
			old, err := utils.ParseCrt(crtName)
			if err != nil {
				return err
			}
			c.DNSNames = old.DNSNames
			c.IPs = old.IPAddresses
			c.CN = old.Subject.CommonName
			c.Organizations = old.Subject.Organization
			c.IsCA = old.IsCA
			c.KeyUsage = old.KeyUsage
			c.ExtKeyUsages = old.ExtKeyUsage
		}
	}
	cert, key, err := g.Gen(c)
	if err != nil {
		return err
	}

	log.Debugf("dumping crt: %s", crtName)
	err = ioutil.WriteFile(crtName, cert, 0400)
	if err != nil {
		return errors.Wrapf(err, "write %s", crtName)
	}
	log.Debugf("dumping key: %s", keyName)
	err = ioutil.WriteFile(keyName, key, 0400)
	if err != nil {
		return errors.Wrapf(err, "write %s", keyName)
	}

	return
}

package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/keys"
)

const (
	KeyType = "RSA PRIVATE KEY"
)

type generator struct {
	bits int
}

func New() keys.Generator {
	return &generator{bits: 2048}
}

func (g *generator) Gen() (crypto.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, g.bits)
	if err != nil {
		return nil, errors.Wrap(err, "rsa keygen")
	}

	return key, nil
}

func (g *generator) Encode(key crypto.Signer) []byte {
	keyPem := &pem.Block{
		Type:  KeyType,
		Bytes: x509.MarshalPKCS1PrivateKey(key.(*rsa.PrivateKey)),
	}

	return pem.EncodeToMemory(keyPem)
}

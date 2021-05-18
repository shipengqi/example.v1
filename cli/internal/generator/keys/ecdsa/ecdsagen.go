package ecdsa

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/keys"
)

const (
	KeyType = "EC PRIVATE KEY"
)

type generator struct {
	curve    elliptic.Curve
}

func New() keys.Generator {
	return &generator{curve: elliptic.P256()}
}

func (g *generator) Gen() (crypto.Signer, error) {
	key, err := ecdsa.GenerateKey(g.curve, rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "ecdsa keygen")
	}

	return key, nil
}

func (g *generator) Encode(key crypto.Signer) []byte {
	x509Encoded, _ := x509.MarshalECPrivateKey(key.(*ecdsa.PrivateKey))
	keyPem := &pem.Block{
		Type:  KeyType,
		Bytes: x509Encoded,
	}

	return pem.EncodeToMemory(keyPem)
}

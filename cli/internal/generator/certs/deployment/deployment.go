package deployment

import (
	"fmt"
	"github.com/shipengqi/example.v1/cli/pkg/kube"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type generator struct {
	namespace string
	vault     *vault.Client
	kube      kube.Client
}

func New() certs.Generator {
	return &generator{}
}

func (g *generator) Gen(c *certs.Certificate) (cert, key []byte, err error) {
	ttl := fmt.Sprintf("%dh", c.Period*24)
	if c.UintTime == types.CertUnitTimeMinute {
		ttl = fmt.Sprintf("%dm", c.Period)
	}

	log.Debugf("generate external certificates for host: %s", c.CN)
	data, err := g.vault.GenerateCert(ttl, c.CN, types.VaultPkiPathRE)
	if err != nil {
		return nil, nil, errors.Wrap(err, "issue cert")
	}
	cert = []byte(data.Certificate + "\n" + data.IssuingCa)
	key = []byte(data.PrivateKey)
	return
}

func (g *generator) Dump(certName, keyName, secret string, cert, key []byte) error {
	data := make(map[string][]byte)
	data[certName] = cert
	data[keyName] = key

	_, err := g.kube.ApplySecretBytes(g.namespace, secret, data)
	if err != nil {
		return err
	}
	return nil
}

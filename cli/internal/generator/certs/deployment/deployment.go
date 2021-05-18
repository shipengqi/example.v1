package deployment

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type generator struct {
	vault *vault.Client
}

func New() certs.Generator {
	return &generator{}
}

func (g *generator) Gen(c *certs.Certificate) (cert, key string, err error) {
	ttl := fmt.Sprintf("%dh", c.Period*24)
	if c.UintTime == certs.CertUnitTimeMinute {
		ttl = fmt.Sprintf("%dm", c.Period)
	}

	log.Debugf("generate external certificates for host: %s", c.CN)
	data, err := g.vault.GenerateCert(ttl, c.CN, certs.VaultREPkiPath)
	if err != nil {
		return "", "", errors.Wrap(err, "issue cert")
	}
	cert = data.Certificate + "\n" + data.IssuingCa
	key = data.PrivateKey
	return
}

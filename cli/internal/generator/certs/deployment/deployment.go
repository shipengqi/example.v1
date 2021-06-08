package deployment

import (
	"encoding/base64"
	"fmt"
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type generator struct {
	namespace string
	vault     *vault.Client
	kube      *kube.Client
}

const (
	SecretNameVaultPass     = "vault-passphrase"
	SecretNameVaultCred     = "vault-credential"
)

func New(ns string, k *kube.Config, v *vault.Config) (certs.Generator, error) {
	kclient, err := kube.New(k)
	if err != nil {
		return nil, err
	}
	vclient, err := vault.New(v)
	if err != nil {
		return nil, err
	}
	return &generator{
		namespace: ns,
		vault:     vclient,
		kube:      kclient,
	}, nil
}

func (g *generator) Gen(c *certs.Certificate) (cert, key []byte, err error) {
	ttl := fmt.Sprintf("%dh", c.Validity*24)
	if c.UintTime == types.CertUnitTimeMinute {
		ttl = fmt.Sprintf("%dm", c.Validity)
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

func (g *generator) GenAndDump(c *certs.Certificate, resources string) (err error) {
	secrets, field, err := parseResources(resources)
	if err != nil {
		return err
	}
	cert, key, err := g.Gen(c)
	if err != nil {
		return err
	}
	data := make(map[string][]byte)
	data[field+".crt"] = cert
	data[field+".key"] = key

	for k := range secrets {
		secret := strings.TrimSpace(secrets[k])
		if len(secret) == 0 {
			continue
		}

		_, err = g.kube.ApplySecretBytes(g.namespace, secret, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) setToken() error {
	vaultPassphrase, err := g.kube.GetSecret(SecretNameVaultPass, g.namespace)
	if err != nil {
		return err
	}
	vaultCredential, err := g.kube.GetSecret(SecretNameVaultCred, g.namespace)
	if err != nil {
		return err
	}

	passphrase := vaultPassphrase.Data["passphrase"]
	encryptedToken := vaultCredential.Data["root.token"]
	encryptedStr := base64.StdEncoding.EncodeToString(encryptedToken)
	tokenEncIv := vaultCredential.Data["root.token.enc_iv"]
	tokenEncKey := vaultCredential.Data["root.token.enc_key"]

	token, err := sysc.ParseVaultToken(encryptedStr, string(passphrase), string(tokenEncKey), string(tokenEncIv))
	if err != nil {
		return err
	}
	g.vault.SetToken(token)
	return nil
}

func parseResources(resources string) ([]string, string, error) {
	rs := strings.Split(resources, " ")
	if len(rs) < 2 {
		return nil, "", errors.New("invalid resources")
	}
	names := strings.Split(rs[0], ",")
	field := rs[1]
	return names, field, nil
}

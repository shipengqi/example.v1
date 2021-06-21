package deployment

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/cli/internal/generator/certs"
	"github.com/shipengqi/example.v1/cli/internal/sysc"
	"github.com/shipengqi/example.v1/cli/internal/types"
	"github.com/shipengqi/example.v1/cli/pkg/kube"
	"github.com/shipengqi/example.v1/cli/pkg/log"
	"github.com/shipengqi/example.v1/cli/pkg/vault"
)

type generator struct {
	namespace    string
	cdfNamespace string
	primary      bool
	vault        *vault.Client
	kube         *kube.Client
}

const (
	SecretNameVaultPass    = "vault-passphrase"
	SecretNameVaultCred    = "vault-credential"
	ResourceKeyPassphrase  = "passphrase"
	ResourceKeyRootToken   = "root.token"
	ResourceKeyTokenEncIv  = "root.token.enc_iv"
	ResourceKeyTokenEncKey = "root.token.enc_key"
)

func New(ns, cdfns string, primary bool, k *kube.Config, v *vault.Config) (certs.Generator, error) {
	kclient, err := kube.New(k)
	if err != nil {
		return nil, err
	}
	vclient, err := vault.New(v)
	if err != nil {
		return nil, err
	}
	return &generator{
		namespace:    ns,
		cdfNamespace: cdfns,
		primary:      primary,
		vault:        vclient,
		kube:         kclient,
	}, nil
}

func (g *generator) Gen(c *certs.Certificate) (cert, key []byte, err error) {
	ttl := fmt.Sprintf("%dh", c.Validity*24)
	if c.UintTime == types.CertUnitTimeMinute {
		ttl = fmt.Sprintf("%dm", c.Validity)
	}

	log.Debugf("generate external certificates for host: %s", c.CN)
	err = g.setToken()
	if err != nil {
		return nil, nil, err
	}
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


	ns := g.namespace

	if g.primary && g.namespace != g.cdfNamespace {
		ns = g.cdfNamespace
	}

	for k := range secrets {
		secret := strings.TrimSpace(secrets[k])
		if len(secret) == 0 {
			continue
		}
		log.Infof("Applying secret: %s in %s ...", secret, ns)
		_, err = g.kube.ApplySecretBytes(ns, secret, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) setToken() error {
	ns := g.namespace
	if g.primary {
		ns = g.cdfNamespace
	}

	vaultPassphrase, err := g.kube.GetSecret(ns, SecretNameVaultPass)
	if err != nil {
		return err
	}
	vaultCredential, err := g.kube.GetSecret(ns, SecretNameVaultCred)
	if err != nil {
		return err
	}

	passphrase := vaultPassphrase.Data[ResourceKeyPassphrase]
	encryptedToken := vaultCredential.Data[ResourceKeyRootToken]
	encryptedStr := base64.StdEncoding.EncodeToString(encryptedToken)
	tokenEncIv := vaultCredential.Data[ResourceKeyTokenEncIv]
	tokenEncKey := vaultCredential.Data[ResourceKeyTokenEncKey]

	token, err := sysc.ParseVaultToken(encryptedStr, string(passphrase), string(tokenEncKey), string(tokenEncIv))
	if err != nil {
		return err
	}
	g.vault.SetToken(token)
	return nil
}

func parseResources(resources string) (names []string, field string, err error) {
	rs := strings.Split(resources, " ")
	if len(rs) < 2 {
		return nil, "", errors.New("invalid resources")
	}
	names = strings.Split(rs[0], ",")
	field = rs[1]
	return names, field, nil
}

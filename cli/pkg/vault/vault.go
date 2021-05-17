package vault

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type Client struct {
	client *api.Client
	config *Config
}

type Config struct {
	Address    string
	Role       string
	ServerName string
	CAs        [][]byte
}

type CrtData struct {
	Certificate string `json:"certificate"`
	IssuingCa   string `json:"issuing_ca"`
	PrivateKey  string `json:"private_key"`
}

func New(conf *Config) (*Client, error) {
	if strings.Trim(conf.Role, " ") == "" {
		return nil, errors.New("vault role is required")
	}
	if strings.Trim(conf.Address, " ") == "" {
		return nil, errors.New("vault address is required")
	}

	caCertPool := x509.NewCertPool()
	for k := range conf.CAs {
		caCertPool.AppendCertsFromPEM(conf.CAs[k])
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:    caCertPool,
			ServerName: conf.ServerName,
		},
	}
	c := &http.Client{Transport: tr}
	client, err := api.NewClient(&api.Config{
		Address:    conf.Address,
		HttpClient: c,
	})

	if err != nil {
		return nil, err
	}
	return &Client{client, conf}, nil
}

func (c *Client) Login(roleId, secretId string) error {
	res, err := c.client.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	})
	if err != nil {
		return err
	}
	c.client.SetToken(res.Auth.ClientToken)
	return nil
}

func (c *Client) SetToken(token string) {
	c.client.SetToken(token)
}

func (c *Client) GenerateCert(ttl string, hostname string, pki string) (*CrtData, error) {
	res, err := c.client.Logical().Write(fmt.Sprintf("%s/issue/%s", pki, c.config.Role), map[string]interface{}{
		"ttl":         ttl,
		"common_name": hostname,
	})
	if err != nil {
		return nil, err
	}
	crt := &CrtData{}
	if _, ok := res.Data["certificate"]; !ok {
		return nil, errors.New("certificate nil")
	}
	if _, ok := res.Data["issuing_ca"]; !ok {
		return nil, errors.New("issuing_ca nil")
	}
	if _, ok := res.Data["private_key"]; !ok {
		return nil, errors.New("private_key nil")
	}

	if v, ok := res.Data["certificate"].(string); !ok {
		return nil, errors.New("certificate type error")
	} else {
		crt.Certificate = v
	}

	if v, ok := res.Data["issuing_ca"].(string); !ok {
		return nil, errors.New("issuing_ca type error")
	} else {
		crt.IssuingCa = v
	}

	if v, ok := res.Data["private_key"].(string); !ok {
		return nil, errors.New("private_key type error")
	} else {
		crt.PrivateKey = v
	}

	return crt, nil
}

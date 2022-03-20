package v1

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

func CreateTransport(fpath string) (*http.Transport, error) {
	pool := x509.NewCertPool()
	caCerts, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(caCerts)

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}, nil
}

package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func CheckCrtValidity(cert *x509.Certificate) int {
	available := cert.NotAfter.Sub(time.Now()).Hours()
	if available <= 0 {
		return -1
	}

	return int(available)
}

func CheckCrtStringValidity(cert *x509.Certificate) int {
	available := cert.NotAfter.Sub(time.Now()).Hours()
	if available <= 0 {
		return -1
	}

	return int(available)
}

func ParseCrt(certPath string) (*x509.Certificate, error) {
	certFile, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	certBlock, _ := pem.Decode(certFile)

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "parse %s", certPath)
	}
	return cert, nil
}

func ParseCrtString(certString string) (*x509.Certificate, error) {
	if len(certString) == 0 {
		return nil, nil
	}
	certBlock, _ := pem.Decode([]byte(certString))
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func ParseKey(keyPath string) (crypto.PrivateKey, error) {
	f, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return ParseKeyString(f, false)
}

func ParseKeyString(keystr []byte, isBase64 bool) (crypto.PrivateKey, error)  {
	var err error
	dkeystr := keystr

	if isBase64 {
		dkeystr, err = base64.StdEncoding.DecodeString(string(keystr))
		if err != nil {
			return nil, errors.Wrap(err, "decode key str")
		}
	}
	privateKeyBlock, _ := pem.Decode(dkeystr)
	if privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes); err == nil {
		return privateKey, nil
	}

	if privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes); err == nil {
		switch privateKey.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return privateKey, nil
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping.")
		}
	}

	if privateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes); err == nil {
		return privateKey, nil
	}
	return nil, errors.Wrap(err, "parse key str")
}

func IsIPV4(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	return true
}

func IsExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsEmptyStr(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}

	return false
}

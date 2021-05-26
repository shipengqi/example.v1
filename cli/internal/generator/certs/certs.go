package certs

type Generator interface {
	Gen(c *Certificate) (cert, key []byte, err error)
	Dump(certName, keyName, secret string, cert, key []byte) error
}

package certs

type Generator interface {
	Gen(c *Certificate) (cert, key []byte, err error)
	GenAndDump(c *Certificate, output string) (err error)
}

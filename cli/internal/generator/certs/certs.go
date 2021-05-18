package certs

const (
	VaultREPkiPath     = "RE"
	CertUnitTimeDay    = "d"
	CertUnitTimeMinute = "m"
)

type Generator interface {
	Gen(c *Certificate) (cert, key string, err error)
}

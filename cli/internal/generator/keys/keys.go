package keys

import "crypto"

type Generator interface {
	Gen() (crypto.Signer, error)
	Encode(key crypto.Signer) []byte
}

package utils

import (
	"encoding/hex"
)

func EncodeXOR(str, key string) string {
	var encrypted []byte
	sl := len(str)
	kl := len(key)

	for i := 0; i < sl; i++ {
		encrypted = append(encrypted, (str[i]) ^ (key[i%kl]))
	}
	return hex.EncodeToString(encrypted)
}

func DecodeXOR(str, key string) (string, error) {
	var decrypted []byte
	de, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	dl := len(de)
	kl := len(key)
	for i := 0; i < dl; i++ {
		decrypted = append(decrypted, (de[i]) ^ (key[i%kl]))
	}

	return string(decrypted), nil
}

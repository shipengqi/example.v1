package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	SALT = "llsfhfhhf$jjfklsjn52522@@44ddddsdfsiwotpvbnusf"
)
// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// EncodeMD5WithSalt md5 encrypt with salt
func EncodeMD5WithSalt(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	m.Write([]byte(SALT))

	return hex.EncodeToString(m.Sum(nil))
}
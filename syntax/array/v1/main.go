package main

import (
	"crypto/x509"
	"fmt"
	"strconv"
)

var keyUsageName = [...]string{
	x509.KeyUsageDigitalSignature: "DigitalSignature",
	x509.KeyUsageContentCommitment: "ContentCommitment",
	x509.KeyUsageKeyEncipherment: "KeyEncipherment",
	x509.KeyUsageDataEncipherment: "DataEncipherment",
	x509.KeyUsageKeyAgreement: "KeyAgreement",
	x509.KeyUsageCertSign: "CertSign",
	x509.KeyUsageCRLSign: "CRLSign",
	x509.KeyUsageEncipherOnly: "EncipherOnly",
	x509.KeyUsageDecipherOnly: "DecipherOnly",
}

func main() {
	fmt.Println("length: ", len(keyUsageName))
	for k, v := range keyUsageName {
		fmt.Println("key: ", k, "value: ", v)
	}
	fmt.Println("keyUsageToString: ", keyUsageToString(x509.KeyUsageCertSign))
}

func keyUsageToString(keyUsage x509.KeyUsage) string {
	if 0 < keyUsage && int(keyUsage) < len(keyUsageName) {
		return keyUsageName[keyUsage]
	}
	return strconv.Itoa(int(keyUsage))
}

// length:  257
// key:  0 value:
// key:  1 value:  DigitalSignature
// key:  2 value:  ContentCommitment
// key:  3 value:
// key:  4 value:  KeyEncipherment
// key:  5 value:
// key:  6 value:
// key:  7 value:
// key:  8 value:  DataEncipherment
// key:  9 value:
// ...
// key:  15 value:
// key:  16 value:  KeyAgreement
// key:  17 value:
// ...
// key:  31 value:
// key:  32 value:  CertSign
// key:  33 value:
// ...
// key:  63 value:
// key:  64 value:  CRLSign
// key:  65 value:
// ...
// key:  127 value:
// key:  128 value:  EncipherOnly
// key:  129 value:
// ...
// key:  255 value:
// key:  256 value:  DecipherOnly
// keyUsageToString:  CertSign
//
// Process finished with exit code 0
// 上面例子中的数组长度是 257， 这是因为
// const (
//	KeyUsageDigitalSignature KeyUsage = 1 << iota
//	KeyUsageContentCommitment
//	KeyUsageKeyEncipherment
//	KeyUsageDataEncipherment
//	KeyUsageKeyAgreement
//	KeyUsageCertSign
//	KeyUsageCRLSign
//	KeyUsageEncipherOnly
//	KeyUsageDecipherOnly
// )
// 上面的 iota 的常量都左移一位，也就是 1, 2, 4, 8, 16, 32, 64, 128, 256

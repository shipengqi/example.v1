package main

import (
	"crypto/rand"
	"encoding/base64"
)

// crypto/rand 是为了提供更好的随机性满足密码对随机数的要求，在 linux 上已经有一个实现就是 /dev/urandom，crypto/rand 就是从这个地
// 方读“真随机”数字返回，但性能比较慢
func main() {
	for i := 0; i < 4; i++  {
		// n, _ := rand.Int(rand.Reader, big.NewInt(1<<62))
		n, _ := rand.Prime(rand.Reader, 256)
		println(n.String())
		strs := base64.StdEncoding.EncodeToString(n.Bytes())
		println(strs)
	}
}

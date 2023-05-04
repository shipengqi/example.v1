package v1

import "testing"

func TestCalculateStrBits(t *testing.T) {
	bitsLen := CalculateStrBits("hello")
	t.Log(bitsLen)

	bitsLen = CalculateStrBits("hello world")
	t.Log(bitsLen)

	bitsLen = CalculateStrBits("你好")
	t.Log(bitsLen)

	bitsLen = CalculateStrBits("你好啊")
	t.Log(bitsLen)

	bitsLen = CalculateStrBits("你好啊！")
	t.Log(bitsLen)
}

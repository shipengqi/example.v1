package substring

import (
	"testing"
)

var testString = "Go 语言是 Google 开发的一种静态类型，编译型，并发型，并具有垃圾回收功能的编程语言。简称为 Golang。"
var asciiString = "DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, " +
	"DNS:autorh78vm00.hpeswlab.net, DNS:autorh78vm00, IP Address:172.17.17.1, IP Address:16.155.195.168"
var asciiString2 = "DNS:kubernetes"
var asciiString3 = "IP:16.155.195.168"
var testLength = 20

func TestSliceSubString(t *testing.T) {
	r := SliceSubString(asciiString, 0, testLength)
	t.Log(r)
	r2 := SliceSubString(asciiString2, 3, len(asciiString2))
	t.Log(r2)
	r3 := SliceSubString(asciiString2, 3, len(asciiString2) - 1)
	t.Log(r3)
	r4 := SliceSubString(asciiString2, 4, len(asciiString2))
	t.Log(r4)
	r5 := SliceSubString(asciiString3, 2, len(asciiString3))
	t.Log(r5)
	r6 := SliceSubString(asciiString3, 3, len(asciiString3))
	t.Log(r6)
}

func BenchmarkRuneSubString(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		RuneSubString(testString, 0, testLength)
	}
}

func BenchmarkDecodeRuneInStringSubString(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		DecodeRuneInStringSubString(testString, 0, testLength)
	}
}

func BenchmarkRangeSubString(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		RangeSubString(testString, 0, testLength)
	}
}

func BenchmarkExtRuneIndexInStringSubString(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		ExtRuneIndexInStringSubString(testString, 0, testLength)
	}
}

func BenchmarkExtRuneSubString(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		ExtRuneSubString(testString, 0, testLength)
	}
}

// Output:
// goos: windows
// goarch: amd64
// pkg: github.com/shipengqi/example.v1/syntax/substring
// BenchmarkRuneSubString-8                         1402653               835 ns/op
//
// BenchmarkDecodeRuneInStringSubString-8          10532376               115 ns/op
//
// BenchmarkRangeSubString-8                       16221542                69.8 ns/op
//
// BenchmarkExtRuneIndexInStringSubString-8        17656743                70.5 ns/op
//
// BenchmarkExtRuneSubString-8                     16450930                73.6 ns/op
// PASS

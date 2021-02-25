package substring

import (
	"testing"
)

var testString = "Go 语言是 Google 开发的一种静态类型，编译型，并发型，并具有垃圾回收功能的编程语言。简称为 Golang。"
var testLength = 20

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

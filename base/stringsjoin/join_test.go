package stringsjoin

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// 直接使用运算符 +
// golang 里面的字符串都是不可变的，
// 每次运算都会产生一个新的字符串，所以会产生很多临时的无用的字符串，不仅没有用，还会给 gc 带来额外的负担，所以性能比较差
func BenchmarkAddStringWithOperator(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = hello + "," + world
	}
}

// fmt.Sprintf
// 内部使用 []byte 实现，不像直接运算符这种会产生很多临时的字符串，
// 但是内部的逻辑比较复杂，有很多额外的判断，还用到了 interface，所以性能也不是很好
func BenchmarkAddStringWithSprintf(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s,%s", hello, world)
	}
}

// strings.Join
// Join 会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符串填入，在已有一个数组的情况下，这种效率会很高，
// 但是本来没有，去构造这个数据的代价也不小
func BenchmarkAddStringWithJoin(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		_ = strings.Join([]string{hello, world}, ",")
	}
}

// buffer.WriteString
// 可以当成可变字符使用，对内存的增长也有优化
func BenchmarkAddStringWithBuffer(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < b.N; i++ {
		var buffer bytes.Buffer
		buffer.WriteString(hello)
		buffer.WriteString(",")
		buffer.WriteString(world)
		_ = buffer.String()
	}
}

// strings.Builder
// 为了改进 buffer 拼接的性能，从 go 1.10 版本开始，增加了一个 builder 类型，用于提升字符串拼接的性能。它的使用和 buffer 几乎一样
func BenchmarkStringBuilder(b *testing.B) {
	hello := "hello"
	world := "world"
	for i:=0;i<b.N;i++{
		var b strings.Builder
		b.WriteString(hello)
		b.WriteString(",")
		b.WriteString(world)
		_ = b.String()
	}
}

// $ go test -bench=.
// Output:
// goos: windows
// goarch: amd64
// pkg: github.com/shipengqi/example.v1/syntax/stringsjoin
// BenchmarkAddStringWithOperator-8        46157041                26.7 ns/op
// BenchmarkAddStringWithSprintf-8          6596523               190 ns/op
// BenchmarkAddStringWithJoin-8            25539522                47.7 ns/op
// BenchmarkAddStringWithBuffer-8          17152364                67.1 ns/op
// BenchmarkStringBuilder-8                19680550                63.3 ns/op
// PASS
// ok      github.com/shipengqi/example.v1/syntax/stringsjoin      6.662s

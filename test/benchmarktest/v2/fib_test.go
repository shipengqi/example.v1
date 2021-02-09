package v2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 执行 `go test -bench=.` 输出
// BenchmarkFib1
// BenchmarkFib1-12     	699541624	         1.72 ns/op
// BenchmarkFib2
// BenchmarkFib2-12     	245036128	         4.90 ns/op
// BenchmarkFib3
// BenchmarkFib3-12     	143020576	         8.38 ns/op
// BenchmarkFib10
// BenchmarkFib10-12    	 3726466	       325 ns/op
// BenchmarkFib20
// BenchmarkFib20-12    	   29256	     39889 ns/op
// BenchmarkFib40
// BenchmarkFib40-12    	       2	 602904550 ns/op
// PASS
//
// 默认情况下，每个基准测试最少运行 1 秒。如果基准测试函数返回时，还不到 1 秒钟，`b.N` 的值会按照序列 1,2,5,10,20,50,... 增加，
// 同时再次运行基准测测试函数。
//
// 我们注意到 `BenchmarkFib40` 一共才运行 2 次。为了更精确的结果，我们可以通过 `-benchtime` 标志指定运行时间，从而使它运行更多次。
// 执行 `go test -bench=Fib40 -benchtime=10s` 输出：
// $ go test -bench=Fib40 -benchtime=10s
// goos: windows
// goarch: amd64
// pkg: github.com/shipengqi/example.v1/test/benchmarktest/v2
// BenchmarkFib40-12             19         620676184 ns/op
// PASS
// ok      github.com/shipengqi/example.v1/test/benchmarktest/v2   13.687s

func BenchmarkFib1(b *testing.B)  { benchmarkFib(1, b) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }

func benchmarkFib(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(i)
	}
}

func TestFib(t *testing.T) {
	var cases = []struct {
		in       int // input
		expected int // expected result
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
	}
	for _, c := range cases {
		actual := Fib(c.in)
		assert.Equal(t, c.expected, actual)
	}
}

func TestFib_Parallel(t *testing.T) {
	var cases = []struct {
		name     string
		in       int // input
		expected int // expected result
	}{
		{"1 的 Fib", 1, 1},
		{"2 的 Fib", 2, 1},
		{"3 的 Fib", 3, 2},
		{"4 的 Fib", 4, 3},
		{"5 的 Fib", 5, 5},
		{"6 的 Fib", 6, 8},
		{"7 的 Fib", 7, 13},
	}
	for _, c := range cases {
		cc := c
		t.Run(c.name, func(t *testing.T) {
			t.Log("time:", time.Now())
			t.Parallel()
			time.Sleep(1 * time.Second)
			actual := Fib(cc.in)
			assert.Equal(t, cc.expected, actual)
		})
	}
}
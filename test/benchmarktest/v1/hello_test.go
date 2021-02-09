package v1

import (
	"bytes"
	"html/template"
	"testing"
)

// 基准函数会运行目标代码 b.N 次。在基准执行期间，程序会自动调整 b.N 直到基准测试函数持续足够长的时间
func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hello("pooky")
	}
}

// Output:
// Hello,  pooky
// ...
// BenchmarkHello-12    	  192034	      5656 ns/op
// PASS
// 上面的输出意味着循环执行了 192034 次，每次循环花费 5656 纳秒 (ns)。


// 如果基准测试在循环前需要一些耗时的配置，则可以先重置定时器
func BenchmarkBigLen(b *testing.B) {
	nbig := NewBig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nbig.Len()
	}
}

// BenchmarkBigLen-12    	771690780	         1.56 ns/op
// PASS

// 如果注释掉 b.ResetTimer()，输出：
// BenchmarkBigLen-12    	       1	1000911100 ns/op
// PASS

// RunParallel 方法能够并行地执行给定的基准测试。
// RunParallel会创建出多个 goroutine，并将 b.N 分配给这些 goroutine 执行，其中 goroutine 数量的默认值为 GOMAXPROCS。
// RunParallel 需要传入一个 body 函数，body 函数将在每个 goroutine 中执行，这个函数需要设置所有 goroutine 本地的状态，
// 并迭代直到 pb.Next 返回 false 值为止。
// StartTimer、StopTime 和 ResetTimer 这三个方法都带有全局作用，所以 body 函数不应该调用这些方法； 除此之外，body 函数也不应该调用 Run 方法。
// ReportAllocs 方法用于打开当前基准测试的内存统计功能， 与 go test 使用 -benchmem 标志类似，但 ReportAllocs 只影响那些调用了该函数的基准测试。
func BenchmarkTemplateParallel(b *testing.B) {
	b.ReportAllocs()
	temp := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		// 每个 goroutine 有属于自己的 bytes.Buffer.
		var buf bytes.Buffer
		for pb.Next() {
			// 循环体在所有 goroutine 中总共执行 b.N 次
			buf.Reset()
			temp.Execute(&buf, "World")
		}
	})
}

// Output：
// BenchmarkTemplateParallel-12    	 4281922	       272 ns/op	     272 B/op	       8 allocs/op
// PASS
// 4281922 ：基准测试的迭代总次数 b.N
// 272 ns/op：平均每次迭代所消耗的纳秒数
// 272 B/op：平均每次迭代内存所分配的字节数
// 8 allocs/op：平均每次迭代的内存分配次数

// testing 包中的 BenchmarkResult 类型能为你提供帮助，它保存了基准测试的结果，定义如下：
// type BenchmarkResult struct {
//    N         int           // The number of iterations. 基准测试的迭代总次数，即 b.N
//    T         time.Duration // The total time taken. 基准测试的总耗时
//    Bytes     int64         // Bytes processed in one iteration. 一次迭代处理的字节数，通过 b.SetBytes 设置
//    MemAllocs uint64        // The total number of memory allocations. 内存分配的总次数
//    MemBytes  uint64        // The total number of bytes allocated. 内存分配的总字节数
// }
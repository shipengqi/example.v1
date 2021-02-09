package v1

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pairs = []struct {
	k string
	v string
}{
	{"key1", "value1"},
	{"key2", "value2"},
	{"key3", "value3"},
	{"key4", "value4"},
	{"key5", "value5"},
	{"key6", "value6"},
	{"key7", "value7"},
	{"key8", "value8"},
	{"key9", "value9"},
	{"key10", "value10"},
	{"key11", "value11"},
	{"key12", "value12"},
	{"key13", "value13"},
}

// 默认情况下 go test 在不同的 package 之间是并行执行测试，在每个 package 内部是串行执行测试。
// 如果想要在 package 内部开启并行测试，需要在测试函数中显式执行 t.Parallel()
// Parallel 方法表示当前测试只会与其他带有 Parallel 方法的测试并行进行测试。
// 如果代码能够进行并行测试，在写测试时，尽量加上 Parallel，这样可以测试出一些可能的问题。
func TestWriteToMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		WriteToMap(tt.k, tt.v)
	}
}

func TestReadFromMap(t *testing.T) {
	t.Parallel()
	for _, tt := range pairs {
		actual := ReadFromMap(tt.k)
		assert.Equal(t, tt.v, actual)
	}
}

func BenchmarkReadFromMap(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}
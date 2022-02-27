package v7

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

// performance

func BenchmarkCasbinFullString(b *testing.B) {
	e, err := casbin.NewEnforcer("testdata/model-with-group-dom.conf", "testdata/policy-with-user-group-dom.csv")
	assert.NoError(b, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		dom      string
		expected bool
	}{
		{"user::pooky", "data1", "read", "tenant::dom1", true},
		{"user::pooky", "data1", "write", "tenant::dom1", true},
		{"user::pooky", "data2", "read", "tenant::dom2", true},
		{"user::pooky", "data2", "write", "tenant::dom2", true},
		{"user::alice", "data1", "read", "tenant::dom1", true},
		{"user::alice", "data1", "write", "tenant::dom1", true},
		{"user::alice", "data2", "read", "tenant::dom2", false},
		{"user::alice", "data2", "write", "tenant::dom2", false},
		{"user::bob", "data1", "read", "tenant::dom1", false},
		{"user::bob", "data1", "write", "tenant::dom1", false},
		{"user::bob", "data2", "read", "tenant::dom2", true},
		{"user::bob", "data2", "write", "tenant::dom2", true},
		{"user::ed", "data1", "read", "tenant::dom1", true},
		{"user::ed", "data1", "write", "tenant::dom1", false},
		{"user::ed", "data2", "read", "tenant::dom2", false},
		{"user::ed", "data2", "write", "tenant::dom2", false},
		{"user::ray", "data1", "read", "tenant::dom1", false},
		{"user::ray", "data1", "write", "tenant::dom1", false},
		{"user::ray", "data2", "read", "tenant::dom2", false},
		{"user::ray", "data2", "write", "tenant::dom2", true},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range tests {
			got, err := e.Enforce(v.sub, v.dom, v.obj, v.act)
			assert.NoError(b, err)
			assert.Equal(b, v.expected, got)
		}
	}
}

func BenchmarkCasbinString(b *testing.B) {
	e, err := casbin.NewEnforcer("testdata/model-with-group-dom.conf", "testdata/policy-with-group-dom.csv")
	assert.NoError(b, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		dom      string
		expected bool
	}{
		{"pooky", "data1", "read", "tenant::dom1", true},
		{"pooky", "data1", "write", "tenant::dom1", true},
		{"pooky", "data2", "read", "tenant::dom2", true},
		{"pooky", "data2", "write", "tenant::dom2", true},
		{"alice", "data1", "read", "tenant::dom1", true},
		{"alice", "data1", "write", "tenant::dom1", true},
		{"alice", "data2", "read", "tenant::dom2", false},
		{"alice", "data2", "write", "tenant::dom2", false},
		{"bob", "data1", "read", "tenant::dom1", false},
		{"bob", "data1", "write", "tenant::dom1", false},
		{"bob", "data2", "read", "tenant::dom2", true},
		{"bob", "data2", "write", "tenant::dom2", true},
		{"ed", "data1", "read", "tenant::dom1", true},
		{"ed", "data1", "write", "tenant::dom1", false},
		{"ed", "data2", "read", "tenant::dom2", false},
		{"ed", "data2", "write", "tenant::dom2", false},
		{"ray", "data1", "read", "tenant::dom1", false},
		{"ray", "data1", "write", "tenant::dom1", false},
		{"ray", "data2", "read", "tenant::dom2", false},
		{"ray", "data2", "write", "tenant::dom2", true},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range tests {
			got, err := e.Enforce(v.sub, v.dom, v.obj, v.act)
			assert.NoError(b, err)
			assert.Equal(b, v.expected, got)
		}
	}
}


func BenchmarkCasbinShortString(b *testing.B) {
	e, err := casbin.NewEnforcer("testdata/model-with-group-dom.conf", "testdata/policy-with-short-group-dom.csv")
	assert.NoError(b, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		dom      string
		expected bool
	}{
		{"pooky", "data1", "read", "tenant::dom1", true},
		{"pooky", "data1", "write", "tenant::dom1", true},
		{"pooky", "data2", "read", "tenant::dom2", true},
		{"pooky", "data2", "write", "tenant::dom2", true},
		{"alice", "data1", "read", "tenant::dom1", true},
		{"alice", "data1", "write", "tenant::dom1", true},
		{"alice", "data2", "read", "tenant::dom2", false},
		{"alice", "data2", "write", "tenant::dom2", false},
		{"bob", "data1", "read", "tenant::dom1", false},
		{"bob", "data1", "write", "tenant::dom1", false},
		{"bob", "data2", "read", "tenant::dom2", true},
		{"bob", "data2", "write", "tenant::dom2", true},
		{"ed", "data1", "read", "tenant::dom1", true},
		{"ed", "data1", "write", "tenant::dom1", false},
		{"ed", "data2", "read", "tenant::dom2", false},
		{"ed", "data2", "write", "tenant::dom2", false},
		{"ray", "data1", "read", "tenant::dom1", false},
		{"ray", "data1", "write", "tenant::dom1", false},
		{"ray", "data2", "read", "tenant::dom2", false},
		{"ray", "data2", "write", "tenant::dom2", true},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range tests {
			got, err := e.Enforce(v.sub, v.dom, v.obj, v.act)
			assert.NoError(b, err)
			assert.Equal(b, v.expected, got)
		}
	}
}

// $ go test -bench=".*" -benchmem -cpu 4
// goos: windows
// goarch: amd64
// pkg: github.com/shipengqi/example.v1/pkgs/sdk/casbin/v7
// cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
// BenchmarkCasbinFullString-4         1940            603588 ns/op          372260 B/op       5845 allocs/op
// BenchmarkCasbinString-4             1980            603396 ns/op          372603 B/op       5845 allocs/op
// BenchmarkCasbinShortString-4        1978            608789 ns/op          371850 B/op       5845 allocs/op
// PASS
// ok      github.com/shipengqi/example.v1/pkgs/sdk/casbin/v7      3.887s


// $ go test -bench=".*" -benchmem -cpu 4 -benchtime 10s
// goos: windows
// goarch: amd64
// pkg: github.com/shipengqi/example.v1/pkgs/sdk/casbin/v7
// cpu: Intel(R) Core(TM) i7-10850H CPU @ 2.70GHz
// BenchmarkCasbinFullString-4        20011            600545 ns/op          372777 B/op       5845 allocs/op
// BenchmarkCasbinString-4            20032            596981 ns/op          372531 B/op       5845 allocs/op
// BenchmarkCasbinShortString-4       19861            601922 ns/op          371651 B/op       5845 allocs/op
// PASS
// ok      github.com/shipengqi/example.v1/pkgs/sdk/casbin/v7      54.28

// benchtime，指定测试时间和循环执行次数（格式需要为Nx，例如100x）
// $ go test -bench=".*" -benchtime=100x -benchmem

// Best model:
// user=username
// group=group::groupname
// role=role::rolename

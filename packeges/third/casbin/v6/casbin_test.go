package v6

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

// 超级管理员

func TestCasbinRBACSimpleRoot(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/simple.conf", "testdata/simple.csv")
	assert.NoError(t, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		dom      string
		expected bool
	}{
		{"alice", "data1", "write", "tenant1", true},
		{"alice", "data1", "read", "tenant1", true},
		{"alice", "data2", "write", "tenant2", true},
		{"alice", "data2", "read", "tenant2", true},
		{"bob", "data1", "read", "tenant1", false},
		{"bob", "data1", "write", "tenant1", false},
		{"bob", "data2", "read", "tenant2", true},
		{"bob", "data2", "write", "tenant2", true},
		{"root", "data1", "write", "tenant1", true},
		{"root", "data1", "read", "tenant1", true},
		{"root", "data2", "write", "tenant2", true},
		{"root", "data2", "read", "tenant2", true},
	}

	for _, v := range tests {
		got, err := e.Enforce(v.sub, v.dom, v.obj, v.act)
		assert.NoError(t, err)
		if v.expected != got {
			t.Logf("enforce failed: %s, %s, %s, %s", v.sub, v.dom, v.obj, v.act)
		}
		assert.Equal(t, v.expected, got)
	}
}

func TestCasbinRBACDomRoot(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/model-with-group-dom.conf", "testdata/policy-with-group-dom.csv")
	assert.NoError(t, err)

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
		{"root", "data1", "read", "tenant::dom1", true},
		{"root", "data1", "write", "tenant::dom1", true},
		{"root", "data2", "read", "tenant::dom2", true},
		{"root", "data2", "write", "tenant::dom2", true},
	}

	for _, v := range tests {
		got, err := e.Enforce(v.sub, v.dom, v.obj, v.act)
		assert.NoError(t, err)
		if v.expected != got {
			t.Logf("enforce failed: %s, %s, %s, %s", v.sub, v.dom, v.obj, v.act)
		}
		assert.Equal(t, v.expected, got)
	}
}

package v5

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

// 优先级模型

func TestCasbinPrioritySimple(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/simple.conf", "testdata/simple.csv")
	assert.NoError(t, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "data1", "write", true}, // for `p, 1, alice, data1, write, allow` has highest priority
		{"alice", "data1", "read", true},
		{"alice", "data2", "write", false},
		{"alice", "data2", "read", false},
		{"bob", "data1", "read", false},
		{"bob", "data1", "write", false},
		{"bob", "data2", "read", false},
		{"bob", "data2", "write", true}, // for bob has role of `data2_allow_group` which has right to write data2, and there's no deny policy with higher priority
	}

	for _, v := range tests {
		got, err := e.Enforce(v.sub, v.obj, v.act)
		assert.NoError(t, err)
		if v.expected != got {
			t.Logf("enforce failed: %s, %s, %s", v.sub, v.obj, v.act)
		}
		assert.Equal(t, v.expected, got)
	}
}

func TestCasbinPriorityRole(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/simple-with-role.conf", "testdata/simple-with-role.csv")
	assert.NoError(t, err)

	tests := []struct {
		sub      string
		obj      string
		act      string
		expected bool
	}{
		{"alice", "data1", "read", true}, // alice 在最底部,所以优先级高于 editor, admin 和 root
		{"jane", "data1", "read", true}, // jane 在最底部,所以优先级高于 editor, admin 和 root
	}

	for _, v := range tests {
		got, err := e.Enforce(v.sub, v.obj, v.act)
		assert.NoError(t, err)
		if v.expected != got {
			t.Logf("enforce failed: %s, %s, %s", v.sub, v.obj, v.act)
		}
		assert.Equal(t, v.expected, got)
	}
}

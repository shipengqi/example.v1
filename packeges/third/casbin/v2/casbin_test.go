package v2

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

func TestCasbinRBAC(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/model-with-group.conf", "testdata/policy-with-group.csv")
	assert.NoError(t, err)

	t.Run("alice authorized", func(t *testing.T) {
		sub := "user::alice" // 想要访问资源的用户。
		obj := "data2" // 将被访问的资源。
		act := "write"  // 用户对资源执行的操作。


		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})
}

func TestCasbinRBACMultiM(t *testing.T) {
	e, err := casbin.NewEnforcer("testdata/model-with-group-multi-m.conf", "testdata/policy-with-group-multi-m.csv")
	assert.NoError(t, err)

	t.Run("alice authorized", func(t *testing.T) {
		sub := "user::alice" // 想要访问资源的用户。
		obj := "data2" // 将被访问的资源。
		act := "write"  // 用户对资源执行的操作。


		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})
}

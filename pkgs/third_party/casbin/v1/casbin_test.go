package v1

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/stretchr/testify/assert"
)

func TestCasbinFileAdapter(t *testing.T) {
	// 新建一个 Casbin 执行器
	// 需要一个Model和一个Adapter
	// 这里使用默认的 FileAdapter
	e, err := casbin.NewEnforcer("testdata/model.conf", "testdata/policy.csv")
	assert.NoError(t, err)

	t.Run("alice authorized", func(t *testing.T) {
		sub := "alice" // 想要访问资源的用户。
		obj := "data1" // 将被访问的资源。
		act := "read"  // 用户对资源执行的操作。

		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("alice authorized failed", func(t *testing.T) {
		sub := "alice" // 想要访问资源的用户。
		obj := "data1" // 将被访问的资源。
		act := "write" // 用户对资源执行的操作。

		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, false, got)

	})
}

func TestCasbinNewModelFromString(t *testing.T) {
	// https://casbin.org/docs/zh-CN/function
	// keyMatch
	// 参数1: 一个 URL 路径，例如 /alice_data/resource1
	// 参数2: 一个 URL 路径，例如 /alice_data/resource1	一个URL 路径或 * 模式下，例如 /alice_data/*

	// keyMatch2
	// 参数1: 一个 URL 路径，例如 /alice_data/resource1
	// 参数2: 一个 URL 路径或 : 模式下，例如 /alice_data/:resource
	var text = `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`
	m, err := model.NewModelFromString(text)
	assert.NoError(t, err)

	// have to use fileadapter.NewAdapter, cannot use directly file path, otherwise it will fail
	a := fileadapter.NewAdapter("testdata/policy.csv")
	e, err := casbin.NewEnforcer(m, a)
	assert.NoError(t, err)

	t.Run("alice authorized", func(t *testing.T) {
		sub := "alice" // 想要访问资源的用户。
		obj := "data1" // 将被访问的资源。
		act := "read"  // 用户对资源执行的操作。

		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("alice authorized failed", func(t *testing.T) {
		sub := "alice" // 想要访问资源的用户。
		obj := "data1" // 将被访问的资源。
		act := "write" // 用户对资源执行的操作。

		got, err := e.Enforce(sub, obj, act)
		assert.NoError(t, err)
		assert.Equal(t, false, got)

	})
}

package v8

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

// API

func TestCasbinAPI(t *testing.T) {
	// 从文件中加载信息
	enforcer, err := casbin.NewEnforcer("./testdata/simple.conf", "./testdata/simple.csv")
	assert.NoError(t, err)
	t.Run("Enforce", func(t *testing.T) {
		got, err := enforcer.Enforce("alice", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("EnforceEx", func(t *testing.T) {
		got, reason, err := enforcer.EnforceEx("amber", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, got)
		assert.Equal(t, []string{"admin", "data1", "read"}, reason)
	})

	t.Run("EnforceEx", func(t *testing.T) {
		got, reason, err := enforcer.EnforceEx("amber", "data1", "read")
		assert.NoError(t, err)
		assert.Equal(t, true, got)
		assert.Equal(t, []string{"admin", "data1", "read"}, reason)
	})

	t.Run("Add, Remove, Update", func(t *testing.T) {
		// 添加一条策略，然后使用 HasPolicy() 来确认
		_, err = enforcer.AddPolicy("added_user", "data1", "read")
		assert.NoError(t, err)

		got := enforcer.HasPolicy("added_user", "data1", "read")
		assert.Equal(t, true, got)

		_, err = enforcer.AddPolicy("added_user2", "data1", "read")
		assert.NoError(t, err)

		got = enforcer.HasPolicy("added_user2", "data1", "read")
		assert.Equal(t, true, got)

		// 移除一条策略，然后使用 HasPolicy() 来确认
		_, err = enforcer.RemovePolicy("added_user2", "data1", "read")
		assert.NoError(t, err)
		got = enforcer.HasPolicy("added_user2", "data1", "read")
		assert.Equal(t, false, got)

		_, err = enforcer.UpdatePolicy([]string{"added_user", "data1", "read"}, []string{"added_user", "data1", "write"})
		assert.NoError(t, err)
		got = enforcer.HasPolicy("added_user", "data1", "read")
		assert.Equal(t, false, got)

		got = enforcer.HasPolicy("added_user", "data1", "write")
		assert.Equal(t, true, got)

		subjects := enforcer.GetAllSubjects()
		assert.Equal(t, []string{"admin", "alice", "bob", "added_user"}, subjects)

		subjects = enforcer.GetAllNamedSubjects("p")
		assert.Equal(t, []string{"admin", "alice", "bob", "added_user"}, subjects)

		objs := enforcer.GetAllObjects()
		assert.Equal(t, []string{"data1", "data2"}, objs)

		objs = enforcer.GetAllNamedObjects("p")
		assert.Equal(t, []string{"data1", "data2"}, objs)

		actions := enforcer.GetAllActions()
		assert.Equal(t, []string{"read", "write"}, actions)

		policies := enforcer.GetPolicy()
		assert.Equal(t, [][]string{
			{"admin", "data1", "read"},
			{"admin", "data1", "write"},
			{"admin", "data2", "read"},
			{"admin", "data2", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "write"},
			{"added_user", "data1", "write"},
		}, policies)

		policies = enforcer.GetNamedPolicy("p")
		assert.Equal(t, [][]string{
			{"admin", "data1", "read"},
			{"admin", "data1", "write"},
			{"admin", "data2", "read"},
			{"admin", "data2", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "write"},
			{"added_user", "data1", "write"},
		}, policies)

		policies = enforcer.GetFilteredPolicy(0, "alice")
		assert.Equal(t, [][]string{{"alice", "data1", "read"}}, policies)

		// panic, must be "p"
		// _ ,err = enforcer.AddNamedPolicy("g", "testadd_role", "admin")

		err = enforcer.SavePolicy()
		assert.NoError(t, err)

		// 移除一条策略，然后使用 HasPolicy() 来确认
		_, err = enforcer.RemovePolicy("added_user", "data1", "write")
		assert.NoError(t, err)
		got = enforcer.HasPolicy("added_user", "data1", "write")
		assert.Equal(t, false, got)

		err = enforcer.SavePolicy()
		assert.NoError(t, err)
	})

}

func TestCasbinRBACAPI(t *testing.T) {
	enforcer, err := casbin.NewEnforcer("./testdata/model-with-group.conf", "./testdata/policy-with-group.csv")
	assert.NoError(t, err)

	t.Run("Grouping APIs", func(t *testing.T) {
		got := enforcer.HasGroupingPolicy("alice", "group::data2::admin")
		assert.Equal(t, true, got)

		got = enforcer.HasGroupingPolicy("alice", "role::data2::read")
		assert.Equal(t, false, got)

		got = enforcer.HasGroupingPolicy("group::data2::admin", "role::data2::read")
		assert.Equal(t, true, got)

		got = enforcer.HasNamedGroupingPolicy("g", "group::data2::admin", "role::data2::write")
		assert.Equal(t, true, got)

		_, err := enforcer.AddGroupingPolicy("group::data2::admin2", "group::data2::admin")
		assert.NoError(t, err)

		got = enforcer.HasGroupingPolicy("group::data2::admin2", "group::data2::admin")
		assert.Equal(t, true, got)

		rules := [][]string{
			{"group::data2::admin2", "role::data2::read"},
		}
		_, err = enforcer.AddGroupingPolicies(rules)
		assert.NoError(t, err)

		err = enforcer.SavePolicy()
		assert.NoError(t, err)

		rules = [][]string{
			{"group::data2::admin2", "role::data2::read"},
			{"group::data2::admin2", "group::data2::admin"},
		}
		_, err = enforcer.RemoveGroupingPolicies(rules)
		assert.NoError(t, err)

		err = enforcer.SavePolicy()
		assert.NoError(t, err)

		got = enforcer.HasGroupingPolicy("group::data2::admin2", "group::data2::admin")
		assert.Equal(t, false, got)
	})

	t.Run("RBAC APIs", func(t *testing.T) {
		got, err := enforcer.GetRolesForUser("alice")
		assert.NoError(t, err)
		assert.Equal(t, []string{"group::data2::admin"}, got)

		got, err = enforcer.GetUsersForRole("group::data2::admin")
		assert.NoError(t, err)
		assert.Equal(t, []string{"alice"}, got)

		has, err := enforcer.HasRoleForUser("alice", "group::data2::admin")
		assert.NoError(t, err)
		assert.Equal(t, true, has)

		has, err = enforcer.HasRoleForUser("alice", "role::data2::write")
		assert.NoError(t, err)
		assert.Equal(t, false, has)

		got, err = enforcer.GetImplicitRolesForUser("alice")
		assert.NoError(t, err)
		assert.Equal(t, []string{"group::data2::admin", "role::data2::read", "role::data2::write"}, got)

		got, err = enforcer.GetImplicitUsersForRole("role::data2::write")
		assert.NoError(t, err)
		assert.Equal(t, []string{"group::data2::admin", "alice"}, got)

		ps := enforcer.GetPermissionsForUser("alice")
		assert.Equal(t, [][]string{}, ps)

		has = enforcer.HasPermissionForUser("alice", "read")
		assert.Equal(t, false, has)

		ps, err = enforcer.GetImplicitPermissionsForUser("alice")
		assert.NoError(t, err)
		assert.Equal(t, [][]string{{"role::data2::read", "data2", "read"}, {"role::data2::write", "data2", "write"}}, ps)

		roles := []string{"data1_admin", "data2_admin"}
		_, err = enforcer.AddRolesForUser("alice", roles)
		assert.NoError(t, err)

		pss := []string{"data3", "read"}
		_, err = enforcer.AddPermissionsForUser("alice", pss)
		assert.NoError(t, err)

		err = enforcer.SavePolicy()
		assert.NoError(t, err)

		_, err = enforcer.DeleteRoleForUser("alice", "data1_admin")
		assert.NoError(t, err)

		_, err = enforcer.DeleteRoleForUser("alice", "data2_admin")
		assert.NoError(t, err)

		_, err = enforcer.DeletePermissionForUser("alice", "data3", "read")
		assert.NoError(t, err)

		err = enforcer.SavePolicy()
		assert.NoError(t, err)
	})
}

func TestCasbinRBACDOMAPI(t *testing.T) {
	enforcer, err := casbin.NewEnforcer("./testdata/model-with-group-dom.conf", "./testdata/policy-with-group-dom.csv")
	assert.NoError(t, err)
	t.Run("GetRolesForUserInDomain", func(t *testing.T) {
		getroletests := []struct {
			user     string
			domain   string
			expected []string
		}{
			{"pooky", "tenant1", []string{"group::data::admin"}},
			{"pooky", "tenant2", []string{"group::data::admin"}},
			{"alice", "tenant1", []string{"group::data::admin"}},
			{"alice", "tenant2", []string{}},
			{"bob", "tenant1", []string{}},
			{"bob", "tenant2", []string{"group::data::admin"}},
			{"ed", "tenant1", []string{"group::data::admin::dom1"}},
			{"ed", "tenant2", []string{}},
			{"ray", "tenant1", []string{}},
			{"ray", "tenant2", []string{"group::data::admin::dom2"}},
		}
		for _, v := range getroletests {
			got := enforcer.GetRolesForUserInDomain(v.user, v.domain)
			assert.Equal(t, v.expected, got)
		}
	})

	t.Run("GetRolesForUserInDomain", func(t *testing.T) {
		getroletests := []struct {
			role     string
			domain   string
			expected []string
		}{
			{"role::data::read", "tenant1", []string{"group::data::admin", "group::data::admin::dom1"}},
			{"role::data::read", "tenant2", []string{"group::data::admin"}},
			{"role::data::write", "tenant1", []string{"group::data::admin"}},
			{"role::data::write", "tenant2", []string{"group::data::admin", "group::data::admin::dom2"}},
			{"group::data::admin", "tenant1", []string{"pooky", "alice"}},
			{"group::data::admin", "tenant2", []string{"pooky", "bob"}},
			{"group::data::admin::dom1", "tenant1", []string{"ed"}},
			{"group::data::admin::dom1", "tenant2", []string{}},
			{"group::data::admin::dom2", "tenant1", []string{}},
			{"group::data::admin::dom2", "tenant2", []string{"ray"}},
		}
		for _, v := range getroletests {
			got := enforcer.GetUsersForRoleInDomain(v.role, v.domain)
			assert.ElementsMatch(t, v.expected, got)
		}
	})

	t.Run("GetPermissionsForUserInDomain", func(t *testing.T) {
		tests := []struct {
			user     string
			domain   string
			expected [][]string
		}{
			{"pooky", "tenant1", [][]string{}},
			{"pooky", "tenant2", [][]string{}},
			{"alice", "tenant1", [][]string{}},
			{"alice", "tenant2", [][]string{}},
			{"bob", "tenant1", [][]string{}},
			{"bob", "tenant2", [][]string{}},
			{"ed", "tenant1", [][]string{}},
			{"ed", "tenant2", [][]string{}},
			{"ray", "tenant1", [][]string{}},
			{"ray", "tenant2", [][]string{}},
			{"role::data::read", "tenant1", [][]string{{"role::data::read", "tenant1", "data1", "read"}}},
			{"role::data::read", "tenant2", [][]string{{"role::data::read", "tenant2", "data2", "read"}}},
			{"role::data::write", "tenant1", [][]string{{"role::data::write", "tenant1", "data1", "write"}}},
			{"role::data::write", "tenant2", [][]string{{"role::data::write", "tenant2", "data2", "write"}}},
		}
		for _, v := range tests {
			got := enforcer.GetPermissionsForUserInDomain(v.user, v.domain)
			assert.ElementsMatch(t, v.expected, got)
		}
	})
}

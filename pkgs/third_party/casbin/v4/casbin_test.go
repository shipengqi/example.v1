package v4

import (
	"fmt"
	"testing"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
)

// ABAC

// API
// 定义用户环境和对象
type Person struct {
	Role string
	Name string
}

type Gate struct {
	Name string
}

type Env struct {
	Time     time.Time
	Location string
}

func (env *Env) IsSchoolTime() bool {
	return env.Time.Hour() >= 8 && env.Time.Hour() <= 18
}

func InitEnv(hour int) *Env {
	env := &Env{}
	env.Time = time.Date(2019, 8, 20, hour, 0, 0, 0, time.Local)
	return env
}

func TestCasbinABAC(t *testing.T) {
	p1 := Person{Role: "Student", Name: "Yun"}
	p2 := Person{Role: "Teacher", Name: "Devin"}
	persons := []Person{p1, p2}
	g1 := Gate{Name: "School Gate"}
	g2 := Gate{Name: "Factory Gate"}
	gates := []Gate{g1, g2}

	t.Run("Basic", func(t *testing.T) {
		// 从文件中加载信息
		enforcer, err := casbin.NewEnforcer("./testdata/abac.conf")
		assert.NoError(t, err)
		envs := []*Env{InitEnv(9), InitEnv(23)}
		for _, env := range envs {
			fmt.Println("\r\nTime:", env.Time.Local())
			for _, p := range persons {
				for _, g := range gates {
					pass, err := enforcer.Enforce(p, g, "In", env)
					assert.NoError(t, err)
					fmt.Println(p.Role, p.Name, "In", g.Name, pass)
					pass, err = enforcer.Enforce(p, g, "Control", env)
					assert.NoError(t, err)
					fmt.Println(p.Role, p.Name, "Control", g.Name, pass)
				}
			}
		}
	})

	t.Run("IsSchoolTime", func(t *testing.T) {
		// 从文件中加载信息
		enforcer, err := casbin.NewEnforcer("./testdata/abac2.conf")
		assert.NoError(t, err)
		envs := []*Env{InitEnv(9), InitEnv(23)}
		for _, env := range envs {
			fmt.Println("\r\nTime:", env.Time.Local())
			for _, p := range persons {
				for _, g := range gates {
					pass, err := enforcer.Enforce(p, g, "In", env)
					assert.NoError(t, err)
					fmt.Println(p.Role, p.Name, "In", g.Name, pass)
					pass, err = enforcer.Enforce(p, g, "Control", env)
					assert.NoError(t, err)
					fmt.Println(p.Role, p.Name, "Control", g.Name, pass)
				}
			}
		}
	})

	t.Run("Complex Rules", func(t *testing.T) {
		// 从文件中加载信息
		enforcer, err := casbin.NewEnforcer("./testdata/model-complex.conf", "./testdata/policy-complex.csv")
		assert.NoError(t, err)
		tests := []struct {
			sub      Person
			obj      Gate
			env      *Env
			act      string
			expected bool
		}{
			{
				Person{
					Role: "Teacher",
					Name: "Lee",
				},
				Gate{Name: "School Gate"},
				InitEnv(9), "In", true,
			},
			{
				Person{
					Role: "Teacher",
					Name: "Lee",
				},
				Gate{Name: "School Gate"},
				InitEnv(23), "Out", true,
			},
			{
				Person{
					Role: "Student",
					Name: "Ming",
				},
				Gate{Name: "School Gate"},
				InitEnv(9), "In", true,
			},
			{
				Person{
					Role: "Student",
					Name: "Ming",
				},
				Gate{Name: "School Gate"},
				InitEnv(9), "Out", false,
			},
			{
				Person{
					Role: "Student",
					Name: "Ming",
				},
				Gate{Name: "School Gate"},
				InitEnv(23), "In", false,
			},
			{
				Person{
					Role: "Student",
					Name: "Ming",
				},
				Gate{Name: "School Gate"},
				InitEnv(23), "Out", true,
			},
		}
		for _, te := range tests {
			fmt.Println("\r\nTime:", te.env.Time.Local())
			got, err := enforcer.Enforce(te.sub, te.obj.Name, te.act, te.env)
			assert.NoError(t, err)
			fmt.Println(te.sub.Role, te.sub.Name, te.act, te.obj.Name, got)
			assert.Equal(t, te.expected, got)
		}
	})
}

package jsons

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 默认情况下序列化与反序列化使用的都是结构体的原生字段名，通过给结构体字段添加 json tag 来指定序列化后的字段名

type Student struct {
	Name   string `json:"name"`
	Age    int64
	Weight float64
	Height float64
}

// 使用 - 作为相应字段 json tag，表示不进行序列化
// 反序列化时仍将恢复该字段，零值处理。

type StudentIgnore struct {
	Name   string `json:"name"`
	Age    int64
	Weight float64 `json:"-"`
	Height float64
}

// 使用 omitempty 作为相应字段 json tag，注意此时在 omitempty 前一定指定一个字段名，否则 omitempty 将作为字段名处理。
// 序列化的时候忽略0值或者空值

type StudentOmitempty struct {
	Name   string `json:"name"`
	Age    int64
	Weight float64 `json:"weight,omitempty"`
	Height float64
}

// 有时在序列化或者反序列化时，可能结构体类型和需要的类型不一致，这个时候可以指定,支持 string, number 和 boolean

type StudentType struct {
	Name string `json:"name"`
	Age  int64  `json:"age,string"`
}

type Body struct {
	Weight float64
	Height float64
}

// `json:",inline"` 通常作用于内嵌的结构体类型
type StudentBody struct {
	Name  string `json:"name"`
	Age   int64  `json:"age,string"`
	*Body `json:",inline"`
}

// 若要在被嵌套结构体整体为空时使其在序列化结果中被忽略，不仅要在被嵌套结构体字段后加上 json:"fileName,omitempty" 还要将其改为结构体指针
type StudentBodyWithoutInline struct {
	Name  string `json:"name"`
	Age   int64  `json:"age,string"`
	*Body `json:"body,omitempty"`
}

func TestTags(t *testing.T) {
	s1 := Student{
		Name:   "jack",
		Age:    20,
		Weight: 71.5,
		Height: 172.5,
	}
	b, err := json.Marshal(s1)
	assert.NoError(t, err)
	t.Logf("s1: %s\n", b) // s1: {"name":"jack","Age":20,"Weight":71.5,"Height":172.5}

	var s2 Student
	err = json.Unmarshal(b, &s2)
	assert.NoError(t, err)

	t.Logf("s2: %#v\n", s2) // s2: jsons.Student{Name:"jack", Age:20, Weight:71.5, Height:172.5}
}

func TestIgnoreTag(t *testing.T) {
	s1 := StudentIgnore{
		Name:   "jack",
		Age:    20,
		Weight: 71.5,
		Height: 172.5,
	}
	b, err := json.Marshal(s1)
	assert.NoError(t, err)
	t.Logf("s1: %s\n", b) // s1: {"name":"jack","Age":20,"Height":172.5}

	var s2 StudentIgnore
	err = json.Unmarshal(b, &s2)
	assert.NoError(t, err)

	t.Logf("s2: %#v\n", s2) // s2: jsons.StudentIgnore{Name:"jack", Age:20, Weight:0, Height:172.5}
}

func TestOmitemptyTag(t *testing.T) {
	s1 := StudentOmitempty{
		Name:   "jack",
		Age:    20,
		Height: 172.5,
	}
	b, err := json.Marshal(s1)
	assert.NoError(t, err)
	t.Logf("s1: %s\n", b) // s1: {"name":"jack","Age":20,"Height":172.5}

	var s2 StudentOmitempty
	err = json.Unmarshal(b, &s2)
	assert.NoError(t, err)

	t.Logf("s2: %#v\n", s2) // jsons.StudentOmitempty{Name:"jack", Age:20, Weight:0, Height:172.5}
}

func TestTypeTag(t *testing.T) {
	s1 := StudentType{
		Name: "jack",
		Age:  20,
	}
	b, err := json.Marshal(s1)
	assert.NoError(t, err)
	t.Logf("s1: %s\n", b) // s1: {"name":"jack","age":"20"}

	var s2 StudentType
	err = json.Unmarshal(b, &s2)
	assert.NoError(t, err)

	t.Logf("s2: %#v\n", s2) // s2: jsons.StudentType{Name:"jack", Age:20}
}

func TestInlineTag(t *testing.T) {
	t.Run("inline tag", func(t *testing.T) {
		s1 := StudentBody{
			Name: "jack",
			Age:  20,
			Body: &Body{
				Weight: 71.5,
				Height: 172.5,
			},
		}
		b, err := json.Marshal(s1)
		assert.NoError(t, err)
		t.Logf("s1: %s\n", b) // s1: {"name":"jack","age":"20","Weight":71.5,"Height":172.5}

		var s2 StudentBody
		err = json.Unmarshal(b, &s2)
		assert.NoError(t, err)

		t.Logf("s2: %#v\n", s2) // s2: jsons.StudentBody{Name:"jack", Age:20, Body:(*jsons.Body)(0xc00009f730)}
	})

	t.Run("without inline tag", func(t *testing.T) {
		s1 := StudentBodyWithoutInline{
			Name: "jack",
			Age:  20,
			Body: &Body{
				Weight: 71.5,
				Height: 172.5,
			},
		}
		b, err := json.Marshal(s1)
		assert.NoError(t, err)
		t.Logf("s1: %s\n", b) // s1: {"name":"jack","age":"20","body":{"Weight":71.5,"Height":172.5}}

		var s2 StudentBodyWithoutInline
		err = json.Unmarshal(b, &s2)
		assert.NoError(t, err)

		t.Logf("s2: %#v\n", s2) // s2: jsons.StudentBodyWithoutInline{Name:"jack", Age:20, Body:(*jsons.Body)(0xc00009fab0)}
	})

}

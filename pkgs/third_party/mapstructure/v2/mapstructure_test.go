package v2

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// 嵌套的结构被认为是拥有该结构体名字的另一个字段。下面的两种方式是一样的
type Person struct {
	Name string
}

// 方式一
type Friend struct {
	Person
}

// 方式二
type Friend2 struct {
	Person Person
}

// 将子结构体的字段提到父结构中
type Friend3 struct {
	Person `mapstructure:",squash"`
}

func ExampleEmbededMapping() {
	// 为了正确解码，Person 结构的数据要在 person 键下：
	data := map[string]interface{}{
		"person": map[string]interface{}{"name": "dj"},
	}

	var f Friend
	mapstructure.Decode(data, &f)
	fmt.Println(f.Person.Name)

	var f2 Friend2
	mapstructure.Decode(data, &f2)
	fmt.Println(f2.Person.Name)

	// 不需要嵌套 person 键
	data2 := map[string]interface{}{
		"name": "dj",
	}
	var f3 Friend3
	mapstructure.Decode(data2, &f3)
	fmt.Println(f3.Person.Name)

	// Output:
	// dj
	// dj
	// dj
}

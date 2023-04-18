package v6

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// 如果不想对结构体字段类型和 `map[string]interface{}` 的对应键值做强类型一致的校验。
// 这时可以使用 `WeakDecode/WeakDecodeMetadata` 方法，会尝试做类型转换
// 如果类型转换失败了，会返回错误。

type Person struct {
	Name   string
	Age    int
	Emails []string
}

func ExampleWeakMapping() {
	// name 对应的值 123 是 int 类型，但是在 WeakDecode 中会将其转换为 string 类型以匹配 Person.Name 字段的 string 类型
	// age 的值 "18" 是 string 类型，在 WeakDecode 中会将其转换为 int 类型以匹配 Person.Age 字段的 int 类型。
	m := map[string]interface{}{
		"name":   123,
		"age":    "18",
		"emails": []int{1, 2, 3},
	}

	var p Person
	err := mapstructure.WeakDecode(m, &p)
	if err == nil {
		fmt.Println("person:", p)
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	// person: {123 18 [1 2 3]}
}

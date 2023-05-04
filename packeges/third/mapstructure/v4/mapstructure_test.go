package v4

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// 将 Go 结构体反向解码为 `map[string]interface{}`
// 反向解码时，可以为某些字段设置 `mapstructure:",omitempty"`。这样当这些字段为零值时，就不会出现在结构的 `map[string]interface{}` 中

type Person struct {
	Name string
	Age  int
	Job  string `mapstructure:",omitempty"`
}

func ExampleMappingOmitEmpty() {
	p := &Person{
		Name: "dj",
		Age:  18,
	}

	var m map[string]interface{}
	mapstructure.Decode(p, &m)

	data, _ := json.Marshal(m)
	fmt.Println(string(data))

	// Output:
	// {"Age":18,"Name":"dj"}
}

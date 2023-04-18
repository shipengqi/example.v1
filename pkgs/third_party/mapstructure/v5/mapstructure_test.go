package v5

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// 解码时会产生一些有用的信息，`mapstructure` 可以使用 `Metadata` 收集这些信息。
// mapstructure.go
// type Metadata struct {
//  // 解码成功的键名
// 	Keys   []string
//  // 在源数据中存在，但是目标结构中不存在的键名。
// 	Unused []string
// }

type Person struct {
	Name string
	Age  int
}

func ExampleMappingMetadata() {
	m := map[string]interface{}{
		"name": "dj",
		"age":  18,
		"job":  "programmer",
	}

	var p Person
	var metadata mapstructure.Metadata
	mapstructure.DecodeMetadata(m, &p, &metadata)

	fmt.Printf("keys:%#v unused:%#v\n", metadata.Keys, metadata.Unused)

	// Output:
	// keys:[]string{"Name", "Age"} unused:[]string{"job"}
}

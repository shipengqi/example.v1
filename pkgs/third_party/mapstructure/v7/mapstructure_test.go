package v7

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

// DecoderConfig 可以配置解码器
// mapstructure.go
// type DecoderConfig struct {
//     // ErrorUnused 为 true 时，如果输入中的键值没有与之对应的字段就返回错误；
//     ErrorUnused       bool
//     // ZeroFields 为 true 时，在 Decode 前清空目标 map。为 false 时，则执行的是 map 的合并。用在 struct 到 map 的转换中；
//     ZeroFields        bool
//     // WeaklyTypedInput 实现 WeakDecode/WeakDecodeMetadata 的功能；
//     WeaklyTypedInput  bool
//     // Metadata 不为 nil 时，收集 Metadata 数据；
//     Metadata          *Metadata
//     // Result 为结果对象，在 map 到 struct 的转换中，Result 为 struct 类型。在 struct 到 map 的转换中，Result 为 map 类型；
//     Result            interface{}
//     // TagName 默认使用 mapstructure 作为结构体的标签名，可以通过该字段设置。
//     TagName           string
// }

type Person struct {
	Name string
	Age  int
}

func ExampleDecoder() {
	m := map[string]interface{}{
		"name": 123,
		"age":  "18",
		"job":  "programmer",
	}

	var p Person
	var metadata mapstructure.Metadata

	// 与 WeakDecode 的功能类似
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &p,
		Metadata:         &metadata,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = decoder.Decode(m)
	if err == nil {
		fmt.Println("person:", p)
		fmt.Printf("keys:%#v, unused:%#v\n", metadata.Keys, metadata.Unused)
	} else {
		fmt.Println(err.Error())
	}

	// Output:
	// person: {123 18}
	// keys:[]string{"Name", "Age"}, unused:[]string{"job"}
}

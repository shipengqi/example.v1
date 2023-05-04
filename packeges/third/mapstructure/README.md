# mapstructure

[mapstructure](github.com/mitchellh/mapstructure) 用于将通用的 `map[string]interface{}` 解码到对应的 Go 结构体中，或者执行相反的操作。


`mapstructure` 使用结构体中字段的名称做这个映射，例如我们的结构体有一个 `Name` 字段，`mapstructure` 解码时会在 `map[string]interface{}` 中查找键名 `name`。

`mapstructure` 处理字段映射是大小写不敏感的。

也可以使用标签指定映射的字段名。

```go
type Person struct {
  Name string `mapstructure:"username"`
}
```

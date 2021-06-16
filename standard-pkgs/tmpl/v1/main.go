package main

import (
	"os"
	"text/template"
)

func main() {
	// 模板定义
	tepl := "My name is {{ . }}"

	// 解析模板
	tmpl, err := template.New("demo").Parse(tepl)

	if err != nil {
		panic(err)
	}
	// 数据驱动模板
	data := "jack"
	// Execute applies a parsed template to the specified data object,
	// and writes the output to wr.
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"os"
	"text/template"
)

type T struct {
	Add func(int) int
}

func (t *T) Sub(i int) int {
	return i - 1
}

func main() {
	ts := &T{
		Add: func(i int) int {
			return i + 1
		},
	}
	tpl := `
		// 只能使用 call 调用
		call field func Add: {{ call .ts.Add .y }}
		// 直接传入 .y 调用
		call method func Sub: {{ .ts.Sub .y }}
	`
	t, _ := template.New("test").Parse(tpl)
	t.Execute(os.Stdout, map[string]interface{}{
		"y": 3,
		"ts": ts,
	})
}

// Output:
//
// 		// 只能使用 call 调用
//		call field func Add: 4
//		// 直接传入 .y 调用
//		call method func Sub: 2


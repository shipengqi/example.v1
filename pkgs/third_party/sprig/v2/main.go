package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
)

const (
	layout = `Test: 
{{wrapWith .Width "\n\t" .Name }}
{{- if .Warnings}}
Warnings:
{{- range .Warnings}}
    {{. | red}}
{{- end}}
{{- end}}
{{ colorize .Name "red" }}
{{ colorize .Name "yellow" }}
{{ colorize .Name "unknown" }}`
)

var redColor = color.New(color.Bold, color.FgRed)

type ExampleT struct {
	Width    int
	Name     string
	Warnings []string
}

func main() {
	funcmap := sprig.TxtFuncMap()
	extras := template.FuncMap{
		"red":      red,
		"colorize": colorize,
	}
	for k, v := range extras {
		funcmap[k] = v
	}
	t := template.New("example template").Funcs(funcmap)
	t2, err := t.Parse(layout)
	if err != nil {
		log.Println(err)
		return
	}
	var buffer bytes.Buffer
	w := bufio.NewWriter(&buffer)
	err = t2.Execute(w, ExampleT{
		Name: "Hello, world",
		Warnings: []string{
			"This is warning1",
			"This is warning2",
		},
		Width: 80,
	})
	if err != nil {
		log.Println(err)
		return
	}
	_ = w.Flush()
	fmt.Println(buffer.String())
	// Test:
	// Hello, world
	// Warnings:
	//    This is warning1
	//    This is warning2
	// Hello, world
	// Hello, world
	// Hello, world
}

func red(text string) string {
	return redColor.SprintfFunc()("%s", text)
}

func colorize(text, col string) string {
	switch strings.ToLower(col) {
	case "red":
		return color.HiRedString(text)
	case "yellow":
		return color.HiYellowString(text)
	default:
		return color.HiGreenString(text)
	}
}

package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	white = "15"
	red   = "9"
	green = "10"
)

const (
	CharPoint   = "•"
	CharLPoint  = "●"
	CharSquare  = "▪"
	CharRhombus = "♦"
)

type Validations struct {
	items  []Validation
	indent string
	char   string
}

// NewValidations Create a new Validations with default options.
func NewValidations(validations ...Validation) *Validations {
	if len(validations) == 0 {
		return nil
	}
	return &Validations{
		items:  validations,
		indent: "",
		char:   CharSquare,
	}
}

func (v Validations) Validate(value string) bool {
	for _, item := range v.items {
		if item.fn != nil && !item.fn(value) {
			return false
		}
	}
	return true
}

func (v Validations) String(value string) string {
	if len(v.items) == 0 {
		return ""
	}
	col := white
	var b strings.Builder
	for _, item := range v.items {
		b.WriteString(v.indent)
		col = green
		if item.fn != nil && !item.fn(value) {
			col = red
		}
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render(v.char))
		b.WriteString(" ")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(item.msg))
		b.WriteString("\n")
	}
	return b.String()
}

type Validation struct {
	fn  func(value string) bool
	msg string
}

func NewValidation(fn func(value string) bool, msg string) Validation {
	return Validation{
		fn:  fn,
		msg: msg,
	}
}

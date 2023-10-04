package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	white = "15"
	red   = "9"
	green = "10"
)

func main() {
	p := tea.NewProgram(NewInputModel("What's your name?", "must larger than 1"))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type validations []validation

func (v validations) Run(value string) string {
	if len(v) == 0 {
		return ""
	}
	col := white
	lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render("•")
	var b strings.Builder
	for _, item := range v {
		// b.WriteString("  ")
		if item.fn != nil {
			col = green
			if err := item.fn(value); err != nil {
				col = red
			}
		}

		// ● • ▪ ♦
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(col)).Render("▪"))
		b.WriteString(" ")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(item.msg))
		b.WriteString("\n")
	}
	return b.String()
}

type validation struct {
	fn  func(value string) error
	msg string
}

type InputModel struct {
	ti          textinput.Model
	validations validations
	label       string
	err         error
}

func NewInputModel(label, placeholder string) InputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return InputModel{
		ti:    ti,
		label: label,
		err:   nil,
		validations: validations{
			validation{
				fn: func(value string) error {
					if len(value) < 2 {
						return errors.New("must large than 1")
					}
					return nil
				},
				msg: "must large than 1",
			},
			validation{
				fn: func(value string) error {
					if len(value) < 1 {
						return errors.New("must less than 5")
					}
					if len(value) > 5 {
						return errors.New("must less than 5")
					}
					return nil
				},
				msg: "must less than 5",
			},
		},
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.ti, cmd = m.ti.Update(msg)
	m.ti.Value()
	return m, cmd
}

func (m InputModel) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s%s",
		m.label,
		m.ti.View(),
		m.validations.Run(m.ti.Value()),
		"(esc to quit)",
	) + "\n"
}

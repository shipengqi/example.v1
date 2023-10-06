package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextAreaModel struct {
	textarea textarea.Model
	label    string
}

func NewTextAreaModel(label string, height int) TextAreaModel {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.SetHeight(height)
	ti.Focus()

	return TextAreaModel{
		textarea: ti,
		label:    label,
	}
}

func (m TextAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TextAreaModel) Value() string {
	return m.textarea.Value()
}

func (m TextAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m TextAreaModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.label,
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}

func main() {
	p := tea.NewProgram(NewTextAreaModel("Tell me a story.", 5))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

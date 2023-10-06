package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TextAreaModel struct {
	textarea textarea.Model
	label    string
	done     bool
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
		case tea.KeyCtrlS, tea.KeyCtrlC:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
			m.done = true
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
	if m.done {
		return fmt.Sprintf(
			"%s\n%s\n",
			m.label,
			lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Render(m.Value()),
		)
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.label,
		m.textarea.View(),
		"(ctrl+s to save)",
	) + "\n\n"
}

func main() {
	p := tea.NewProgram(NewTextAreaModel("Tell me a story.", 5))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	p = tea.NewProgram(NewTextAreaModel("Tell me another story.", 5))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

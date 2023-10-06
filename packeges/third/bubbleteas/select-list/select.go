package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SelectModel struct {
	cursor  int
	choice  string
	label   string
	choices []string
	done    bool
}

func NewSelectModel(label string, choices []string) SelectModel {
	return SelectModel{
		label:   label,
		choices: choices,
	}
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Value() string {
	return m.choice
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			m.done = true
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}

	return m, nil
}

func (m SelectModel) View() string {
	if m.done {
		return fmt.Sprintf(
			"%s %s\n",
			m.label,
			lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Render(m.Value()),
		)
	}

	s := strings.Builder{}
	s.WriteString(m.label)
	s.WriteString("\n")
	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render("[x] "))
			s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render(m.choices[i]))
		} else {
			s.WriteString("[ ] ")
			s.WriteString(m.choices[i])
		}

		s.WriteString("\n")
	}

	return s.String()
}

func main() {
	p := tea.NewProgram(NewSelectModel("What to do today?", []string{
		"Plant carrots",
		"Go to the market",
		"Read something",
		"See friends",
	}))

	// Run returns the model as a tea.Model.
	// m, err := p.Run()
	_, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	// if m, ok := m.(SelectModel); ok {
	// 	fmt.Printf("\n---\nYou chose %s!\n", m.Value())
	// }
}

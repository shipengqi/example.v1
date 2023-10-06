package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"log"
	"os"
)

type InputModel struct {
	ti          textinput.Model
	validations *Validations
	label       string
	done        bool
}

func NewInputModel(label, placeholder string) InputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	im := InputModel{
		ti:    ti,
		label: label,
	}
	v1 := NewValidation(func(value string) bool {
		if len(value) < 2 {
			return false
		}
		return true
	}, "must large than 1")
	v2 := NewValidation(func(value string) bool {
		if len(value) < 1 || len(value) > 5 {
			return false
		}
		return true
	}, "must less than 5")
	im.validations = NewValidations(v1, v2)
	return im
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Value() string {
	return m.ti.Value()
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if !m.validations.Validate(m.Value()) {
				return m, m.ti.Focus()
			}
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			quit <- struct{}{}
			m.ti.Blur()
			return m, tea.Quit
		}
	}

	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	if m.done {
		return fmt.Sprintf(
			"%s %s\n",
			m.label,
			lipgloss.NewStyle().Foreground(lipgloss.Color("43")).Render(m.Value()),
		)
	}
	return fmt.Sprintf(
		"%s\n%s\n%s\n",
		m.label,
		m.ti.View(),
		m.validations.String(m.Value()),
	)
}

var quit = make(chan struct{})

func main() {
	go func() {
		for {
			select {
			case <-quit:
				termenv.DefaultOutput().ShowCursor()
				os.Exit(0)
				return
			}
		}
	}()

	p := tea.NewProgram(NewInputModel("What's your name?", ""))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	p = tea.NewProgram(NewInputModel("How old are you?", ""))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	p = tea.NewProgram(NewInputModel("where are you from?", ""))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

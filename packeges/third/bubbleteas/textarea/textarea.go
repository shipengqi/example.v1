package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	choiceStyle   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("241"))
	saveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("43"))
	quitViewStyle = lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))
)

func filter(teaModel tea.Model, msg tea.Msg) tea.Msg {
	if _, ok := msg.(tea.QuitMsg); !ok {
		return msg
	}

	m := teaModel.(TextAreaModel)
	if m.hasChanges {
		return nil
	}

	return msg
}

type keymap struct {
	save key.Binding
	quit key.Binding
}

type TextAreaModel struct {
	textarea    textarea.Model
	help        help.Model
	keymap      keymap
	label       string
	unSavedWarn bool
	quitting    bool
	hasChanges  bool
	done        bool
}

func NewTextAreaModel(label string, height int, enableUnSavedWarn bool) TextAreaModel {
	ti := textarea.New()
	ti.Placeholder = "Once upon a time..."
	ti.SetHeight(height)
	ti.Focus()

	return TextAreaModel{
		textarea:    ti,
		label:       label,
		help:        help.New(),
		unSavedWarn: enableUnSavedWarn,
		keymap: keymap{
			save: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "save"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
	}
}

func (m TextAreaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TextAreaModel) Value() string {
	return m.textarea.Value()
}

func (m TextAreaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m.updatePromptView(msg)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
			m.hasChanges = false
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
			m.quitting = true
			return m, tea.Quit
		case tea.KeyRunes:
			m.hasChanges = true
			fallthrough // execute the next case regardless of whether it is matched or not.
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
	if m.unSavedWarn && m.quitting {
		if m.hasChanges {
			text := lipgloss.JoinHorizontal(lipgloss.Top, "You have unsaved changes. Quit without saving?", choiceStyle.Render("[y/n]"))
			return quitViewStyle.Render(text)
		}
		return "Very important, thank you\n"
	}
	if m.done {
		return fmt.Sprintf(
			"%s\n%s\n",
			m.label,
			saveTextStyle.Render(m.Value()),
		)
	}
	helpView := m.help.ShortHelpView([]key.Binding{
		m.keymap.save,
		m.keymap.quit,
	})
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.label,
		m.textarea.View(),
		helpView,
	) + "\n\n"
}

func (m TextAreaModel) updatePromptView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// For simplicity's sake, we'll treat any key besides "y" as "no"
		if key.Matches(msg, m.keymap.quit) || msg.String() == "y" {
			m.hasChanges = false
			return m, tea.Quit
		}
		m.quitting = false
	}

	return m, nil
}

func main() {
	// WithFilter supplies an event filter that will be invoked before Bubble Tea
	// processes a tea.Msg. The event filter can return any tea.Msg which will then
	// get handled by Bubble Tea instead of the original event. If the event filter
	// returns nil, the event will be ignored and Bubble Tea will not process it.
	p := tea.NewProgram(NewTextAreaModel("Tell me a story.", 5, true), tea.WithFilter(filter))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	p = tea.NewProgram(NewTextAreaModel("Tell me another story.", 5, false))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

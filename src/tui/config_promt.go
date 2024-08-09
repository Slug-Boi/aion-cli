package tui

// Based on this:
// https://elewis.dev/charming-cobras-with-bubbletea-part-1
// https://github.com/charmbracelet/bubbletea/blob/master/examples/result/main.go

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{}
var terminal_msg = ""

type Model struct {
	cursor int
	choice string
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.choice = "quit"
			return m, tea.Quit

		case "enter", " ":
			// Send the choice on the channel and exit.
			m.choice = choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(terminal_msg + "\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("[x] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func RunConfigTea(cmd_choices []string, cmd_msg string) string {
	// Insert cmd variables into Bubbltea variables
	choices = cmd_choices
	terminal_msg = cmd_msg

	p := tea.NewProgram(Model{})

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
			return "quit"
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(Model); ok && m.choice != "" {
		return m.choice
	}
	return "Something went wrong in bubbleTea command structure oh-uh"
}

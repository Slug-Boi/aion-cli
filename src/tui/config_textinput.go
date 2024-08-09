package tui
// Based on https://github.com/charmbracelet/bubbletea/tree/master/examples/textinput

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func StartTextTUI(title, placeholder string) string {
	viewString = title

	p := tea.NewProgram(initialTextModel(placeholder))
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := m.(TextModel); ok && m.textInput.Value() != "" {

		return "formID:" + m.textInput.Value()
	}
	return "reset"
}

type (
	errMsg error
)

type TextModel struct {
	textInput textinput.Model
	err       error
}

func initialTextModel(placeholder string) TextModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return TextModel{
		textInput: ti,
		err:       nil,
	}
}

func (m TextModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.textInput.SetValue("")
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var viewString = ""

func (m TextModel) View() string {
	return fmt.Sprintf(
		viewString+"\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

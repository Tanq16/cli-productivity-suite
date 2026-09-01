package utils

import (
	"errors"
	"strings"

	tea "charm.land/bubbletea/v2"
)

var ErrNoTerminal = errors.New("no interactive terminal")

type selectModel struct {
	label    string
	options  []string
	cursor   int
	canceled bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		m.cursor = max(m.cursor-1, 0)
	case "down", "j":
		m.cursor = min(m.cursor+1, len(m.options)-1)
	case "enter":
		return m, tea.Quit
	case "esc", "ctrl+c":
		m.canceled = true
		return m, tea.Quit
	}
	return m, nil
}

func (m selectModel) View() tea.View {
	var b strings.Builder
	b.WriteString(infoStyle.Render(m.label) + "\n")
	for i, opt := range m.options {
		if i == m.cursor {
			b.WriteString(successStyle.Render("> "+opt) + "\n")
			continue
		}
		b.WriteString("  " + opt + "\n")
	}
	return tea.NewView(b.String())
}

func PromptSelect(label string, options []string) (int, error) {
	if !StdinIsTerminal {
		return -1, ErrNoTerminal
	}
	final, err := tea.NewProgram(selectModel{label: label, options: options}).Run()
	if err != nil {
		return -1, err
	}
	m := final.(selectModel)
	if m.canceled {
		return -1, nil
	}
	return m.cursor, nil
}

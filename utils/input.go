package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

var stdinScanner *bufio.Scanner

func getStdinScanner() *bufio.Scanner {
	if stdinScanner == nil {
		stdinScanner = bufio.NewScanner(os.Stdin)
	}
	return stdinScanner
}

// ReadPipedLine returns one line, or "" when stdin is a terminal or exhausted.
func ReadPipedLine() string {
	fi, err := os.Stdin.Stat()
	if err != nil || fi.Mode()&os.ModeCharDevice != 0 {
		return ""
	}
	if s := getStdinScanner(); s.Scan() {
		return strings.TrimSpace(s.Text())
	}
	return ""
}

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

// PromptSelect returns the chosen 0-based index, or -1 when the user cancels.
func PromptSelect(label string, options []string) (int, error) {
	if GlobalForAIFlag {
		line := ReadPipedLine()
		if line == "" {
			return -1, nil
		}
		n, err := strconv.Atoi(line)
		if err != nil || n < 1 || n > len(options) {
			return -1, fmt.Errorf("expected a number between 1 and %d, got %q", len(options), line)
		}
		return n - 1, nil
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

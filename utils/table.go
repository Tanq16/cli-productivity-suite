package utils

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.ANSIColor(5))

	evenRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(8))

	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.ANSIColor(8))
)

func PrintTable(headers []string, rows [][]string) {
	if GlobalForAIFlag {
		printMarkdownTable(headers, rows)
		return
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(borderStyle).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if row%2 == 0 {
				return evenRowStyle
			}
			return lipgloss.NewStyle()
		})

	PrintGeneric(t.Render())
}

func printMarkdownTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}
	fmt.Println("| " + strings.Join(escapeCells(headers), " | ") + " |")
	seps := make([]string, len(headers))
	for i := range seps {
		seps[i] = "---"
	}
	fmt.Println("| " + strings.Join(seps, " | ") + " |")
	for _, row := range rows {
		fmt.Println("| " + strings.Join(escapeCells(row), " | ") + " |")
	}
}

func escapeCells(cells []string) []string {
	escaped := make([]string, len(cells))
	for i, cell := range cells {
		escaped[i] = strings.ReplaceAll(cell, "|", "\\|")
	}
	return escaped
}

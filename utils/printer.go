package utils

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"
)

var (
	ColorBlue   = lipgloss.ANSIColor(12)
	ColorGreen  = lipgloss.ANSIColor(10)
	ColorRed    = lipgloss.ANSIColor(9)
	ColorYellow = lipgloss.ANSIColor(11)
	ColorMuted  = lipgloss.ANSIColor(7)
	ColorChrome = lipgloss.ANSIColor(8)

	infoStyle    = lipgloss.NewStyle().Foreground(ColorBlue)
	successStyle = lipgloss.NewStyle().Foreground(ColorGreen)
	errorStyle   = lipgloss.NewStyle().Foreground(ColorRed)
	warnStyle    = lipgloss.NewStyle().Foreground(ColorYellow)
)

func PrintInfo(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println(infoStyle.Render("→ " + msg))
}

func PrintSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println(successStyle.Render("✓ " + msg))
}

func PrintError(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Err(err).Msg(msg)
		} else {
			log.Error().Msg(msg)
		}
		return
	}
	lipgloss.Println(errorStyle.Render("✗ " + msg))
}

func PrintFatal(msg string, err error) {
	PrintError(msg, err)
	os.Exit(1)
}

func PrintWarn(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Warn().Err(err).Msg(msg)
		} else {
			log.Warn().Msg(msg)
		}
		return
	}
	lipgloss.Println(warnStyle.Render("! " + msg))
}

func PrintGeneric(msg string) {
	lipgloss.Println(msg)
}

func PrintIndentedGeneric(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println("    " + msg)
}

func PrintIndentedList(items []string) {
	if GlobalDebugFlag {
		log.Info().Msg(strings.Join(items, ", "))
		return
	}
	width := max(TermWidth()-4, 1)
	var line string
	for i, item := range items {
		if i < len(items)-1 {
			item += ","
		}
		switch {
		case line == "":
			line = item
		case lipgloss.Width(line)+1+lipgloss.Width(item) <= width:
			line += " " + item
		default:
			PrintIndentedGeneric(line)
			line = item
		}
	}
	if line != "" {
		PrintIndentedGeneric(line)
	}
}

func PrintRunning(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println(infoStyle.Render("↻ " + msg))
}

func PrintIndentedSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println(successStyle.Render("  ✓ " + msg))
}

func PrintIndentedError(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Err(err).Msg(msg)
		} else {
			log.Error().Msg(msg)
		}
		return
	}
	lipgloss.Println(errorStyle.Render("  ✗ " + msg))
}

func PrintIndentedWarn(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Warn().Err(err).Msg(msg)
		} else {
			log.Warn().Msg(msg)
		}
		return
	}
	lipgloss.Println(warnStyle.Render("  ! " + msg))
}

func PrintIndentedRunning(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
		return
	}
	lipgloss.Println(infoStyle.Render("  ↻ " + msg))
}

func ClearLines(n int) {
	if GlobalDebugFlag || !StdoutIsTerminal {
		return
	}
	for range n {
		fmt.Print("\033[A\033[2K")
	}
}

func ClearPreviousLine() {
	ClearLines(1)
}

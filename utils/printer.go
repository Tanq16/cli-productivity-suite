package utils

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

var (
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(4))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(3))
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(8))
)

func PrintInfo(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[INFO] " + msg)
	} else {
		fmt.Println(infoStyle.Render("[INFO]") + " " + msg)
	}
}

func PrintSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[OK] " + msg)
	} else {
		fmt.Println(successStyle.Render("[OK]") + " " + msg)
	}
}

func PrintError(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Err(err).Msg(msg)
		} else {
			log.Error().Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[ERROR] " + msg)
	} else {
		fmt.Println(errorStyle.Render("[ERROR]") + " " + msg)
	}
}

func PrintFatal(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Err(err).Msg(msg)
		} else {
			log.Error().Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[FATAL] " + msg)
	} else {
		fmt.Println(errorStyle.Render("[FATAL]") + " " + msg)
	}
	os.Exit(1)
}

func PrintWarn(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Warn().Err(err).Msg(msg)
		} else {
			log.Warn().Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[WARN] " + msg)
	} else {
		fmt.Println(warnStyle.Render("[WARN]") + " " + msg)
	}
}

func PrintGeneric(msg string) {
	fmt.Println(msg)
}

func PrintDebug(msg string) {
	if !GlobalDebugFlag {
		return
	}
	log.Debug().Msg(msg)
}

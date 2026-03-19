package utils

import (
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"
)

var (
	ColorBlue   = lipgloss.Color("12") // Bright Blue
	ColorGreen  = lipgloss.Color("10") // Bright Green
	ColorRed    = lipgloss.Color("9")  // Bright Red
	ColorYellow = lipgloss.Color("11") // Bright Yellow

	infoStyle    = lipgloss.NewStyle().Foreground(ColorBlue)
	successStyle = lipgloss.NewStyle().Foreground(ColorGreen)
	errorStyle   = lipgloss.NewStyle().Foreground(ColorRed)
	warnStyle    = lipgloss.NewStyle().Foreground(ColorYellow)
)

func PrintInfo(msg string) {
	if GlobalDebugFlag {
		log.Info().Str("package", "utils").Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[INFO] " + msg)
	} else {
		fmt.Println(infoStyle.Render("→ " + msg))
	}
}

func PrintSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Str("package", "utils").Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[OK] " + msg)
	} else {
		fmt.Println(successStyle.Render("✓ " + msg))
	}
}

func PrintError(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Str("package", "utils").Err(err).Msg(msg)
		} else {
			log.Error().Str("package", "utils").Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[ERROR] " + msg)
	} else {
		fmt.Println(errorStyle.Render("✗ " + msg))
	}
}

func PrintFatal(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Str("package", "utils").Err(err).Msg(msg)
		} else {
			log.Error().Str("package", "utils").Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[ERROR] " + msg)
	} else {
		fmt.Println(errorStyle.Render("✗ " + msg))
	}
	os.Exit(1)
}

func PrintWarn(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Warn().Str("package", "utils").Err(err).Msg(msg)
		} else {
			log.Warn().Str("package", "utils").Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[WARN] " + msg)
	} else {
		fmt.Println(warnStyle.Render("! " + msg))
	}
}

func PrintGeneric(msg string) {
	fmt.Println(msg)
}

func PrintRunning(msg string) {
	if GlobalDebugFlag {
		log.Info().Str("package", "utils").Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[RUNNING] " + msg)
	} else {
		fmt.Println(infoStyle.Render("↻ " + msg))
	}
}

func PrintIndentedSuccess(msg string) {
	if GlobalDebugFlag {
		log.Info().Str("package", "utils").Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[OK] " + msg)
	} else {
		fmt.Println(successStyle.Render("  ✓ " + msg))
	}
}

func PrintIndentedError(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Error().Str("package", "utils").Err(err).Msg(msg)
		} else {
			log.Error().Str("package", "utils").Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[ERROR] " + msg)
	} else {
		fmt.Println(errorStyle.Render("  ✗ " + msg))
	}
}

func PrintIndentedWarn(msg string, err error) {
	if GlobalDebugFlag {
		if err != nil {
			log.Warn().Str("package", "utils").Err(err).Msg(msg)
		} else {
			log.Warn().Str("package", "utils").Msg(msg)
		}
	} else if GlobalForAIFlag {
		fmt.Println("[WARN] " + msg)
	} else {
		fmt.Println(warnStyle.Render("  ! " + msg))
	}
}

func PrintIndentedRunning(msg string) {
	if GlobalDebugFlag {
		log.Info().Str("package", "utils").Msg(msg)
	} else if GlobalForAIFlag {
		fmt.Println("[RUNNING] " + msg)
	} else {
		fmt.Println(infoStyle.Render("  ↻ " + msg))
	}
}

func ClearLines(n int) {
	if GlobalDebugFlag || GlobalForAIFlag {
		return
	}
	for range n {
		fmt.Print("\033[A\033[2K")
	}
}

func ClearPreviousLine() {
	if GlobalDebugFlag || GlobalForAIFlag {
		return
	}
	fmt.Print("\033[A\033[2K")
}

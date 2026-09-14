package utils

import (
	"os"
	"strconv"

	"github.com/charmbracelet/x/term"
)

const (
	minTermWidth     = 24
	defaultTermWidth = 80
)

var GlobalDebugFlag bool

var (
	StdinIsTerminal  = term.IsTerminal(os.Stdin.Fd())
	StdoutIsTerminal = term.IsTerminal(os.Stdout.Fd())
)

func TermWidth() int {
	if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 0 {
		return w
	}
	if n, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && n >= minTermWidth {
		return n
	}
	return defaultTermWidth
}

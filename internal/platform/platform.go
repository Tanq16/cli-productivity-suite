package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type OS int

const (
	Linux OS = iota
	Darwin
)

func (o OS) String() string {
	switch o {
	case Linux:
		return "linux"
	case Darwin:
		return "darwin"
	default:
		return "unknown"
	}
}

type Arch int

const (
	AMD64 Arch = iota
	ARM64
)

func (a Arch) String() string {
	switch a {
	case AMD64:
		return "amd64"
	case ARM64:
		return "arm64"
	default:
		return "unknown"
	}
}

type Platform struct {
	OS      OS
	Arch    Arch
	HomeDir string
}

func Detect() (Platform, error) {
	var p Platform

	switch runtime.GOOS {
	case "linux":
		p.OS = Linux
	case "darwin":
		p.OS = Darwin
	default:
		return p, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	switch runtime.GOARCH {
	case "amd64":
		p.Arch = AMD64
	case "arm64":
		p.Arch = ARM64
	default:
		return p, fmt.Errorf("unsupported arch: %s", runtime.GOARCH)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return p, fmt.Errorf("cannot determine home directory: %w", err)
	}
	p.HomeDir = home

	return p, nil
}

func (p Platform) ShellDir() string {
	return filepath.Join(p.HomeDir, "shell")
}

func (p Platform) ShellExecDir() string {
	return filepath.Join(p.HomeDir, "shell", "executables")
}

func (p Platform) ShellExtDir() string {
	return filepath.Join(p.HomeDir, "shell", "extensions")
}

func (p Platform) ConfigDir() string {
	return filepath.Join(p.HomeDir, ".config", "cps")
}

func (p Platform) StatePath() string {
	return filepath.Join(p.HomeDir, ".config", "cps", "state.json")
}

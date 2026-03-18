package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
)

// ensureSudo runs "sudo -v" with the terminal attached so the user sees
// the password prompt. The refreshed credential cache covers all
// subsequent sudo calls for 5-15 minutes (OS-dependent).
func ensureSudo() error {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo -v failed: %w", err)
	}
	return nil
}

// phaseNeedsSudo returns true if any of the given tool kinds require sudo
// on the current platform. Used by init to pre-prompt once.
func phaseNeedsSudo(p platform.Platform, kinds ...registry.ToolKind) bool {
	for _, k := range kinds {
		switch k {
		case registry.SystemPackage:
			if p.OS == platform.Linux {
				return true
			}
		case registry.CloudCLI:
			// aws-cli needs sudo on Linux; azure-cli on Linux
			if p.OS == platform.Linux {
				return true
			}
		case registry.LanguageRuntime:
			return true // go-sdk uses sudo on both platforms
		}
	}
	return false
}

// toolNeedsSudo returns true if this specific tool requires sudo on the
// current platform. Used by install for single-tool installs.
func toolNeedsSudo(tool registry.Tool, p platform.Platform) bool {
	switch tool.Kind {
	case registry.SystemPackage:
		return p.OS == platform.Linux
	case registry.CloudCLI:
		switch tool.Name {
		case "aws-cli":
			return p.OS == platform.Linux // macOS uses brew, Linux uses sudo install script
		case "azure-cli":
			return p.OS == platform.Linux
		}
	case registry.LanguageRuntime:
		if tool.Name == "go-sdk" {
			return true // sudo on both platforms (installs to /usr/local/go)
		}
	}
	return false
}

package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
)

func EnsureSudo() error {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sudo -v failed: %w", err)
	}
	return nil
}

func StartSudoRefresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cmd := exec.Command("sudo", "-n", "-v")
				cmd.Stdin = os.Stdin
				cmd.Run()
			}
		}
	}()
}

func PhaseNeedsSudo(p platform.Platform, kinds ...registry.ToolKind) bool {
	for _, k := range kinds {
		switch k {
		case registry.SystemPackage:
			if p.OS == platform.Linux {
				return true
			}
		case registry.CloudCLI:
			if p.OS == platform.Linux {
				return true
			}
		case registry.LanguageRuntime:
			return true
		}
	}
	return false
}

func ToolNeedsSudo(tool registry.Tool, p platform.Platform) bool {
	switch tool.Kind {
	case registry.SystemPackage:
		return p.OS == platform.Linux
	case registry.CloudCLI:
		switch tool.Name {
		case "aws-cli":
			return p.OS == platform.Linux
		case "azure-cli":
			return p.OS == platform.Linux
		}
	case registry.LanguageRuntime:
		if tool.Name == "go-sdk" {
			return true
		}
	}
	return false
}

package installer

import (
	"fmt"
	"os/exec"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type SystemPackageInstaller struct{}

func (s *SystemPackageInstaller) Install(tool *registry.Tool, _ platform.Platform, _ *github.Client, st *state.State) Result {
	if len(tool.BrewPkgs) == 0 && len(tool.BrewCasks) == 0 {
		return Result{Tool: tool.Name, Skipped: true}
	}

	if len(tool.BrewPkgs) > 0 {
		args := append([]string{"install"}, tool.BrewPkgs...)
		cmd := exec.Command("brew", args...)
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("brew install failed: %w", err)}
		}
	}

	if len(tool.BrewCasks) > 0 {
		for _, cask := range tool.BrewCasks {
			cmd := exec.Command("brew", "install", "--cask", cask)
			if err := utils.RunCmd(cmd); err != nil {
				return Result{Tool: tool.Name, Err: fmt.Errorf("brew cask install %s failed: %w", cask, err)}
			}
		}
	}

	st.SetToolVersion(tool.Name, "brew-managed")
	return Result{Tool: tool.Name, Version: "brew-managed"}
}

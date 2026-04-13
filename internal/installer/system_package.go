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

func (s *SystemPackageInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	switch p.OS {
	case platform.Linux:
		return s.installApt(tool, st)
	case platform.Darwin:
		return s.installBrew(tool, st)
	default:
		return Result{Tool: tool.Name, Err: fmt.Errorf("unsupported OS")}
	}
}

func (s *SystemPackageInstaller) installApt(tool *registry.Tool, st *state.State) Result {
	if len(tool.AptPkgs) == 0 {
		return Result{Tool: tool.Name, Skipped: true}
	}

	args := append([]string{"apt-get", "install", "-y"}, tool.AptPkgs...)
	cmd := exec.Command("sudo", args...)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("apt install failed: %w", err)}
	}

	st.SetToolVersion(tool.Name, "system-managed")
	return Result{Tool: tool.Name, Version: "system-managed"}
}

func (s *SystemPackageInstaller) installBrew(tool *registry.Tool, st *state.State) Result {
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

	st.SetToolVersion(tool.Name, "system-managed")
	return Result{Tool: tool.Name, Version: "system-managed"}
}

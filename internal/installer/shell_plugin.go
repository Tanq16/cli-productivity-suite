package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type ShellPluginInstaller struct{}

func (s *ShellPluginInstaller) Check(tool *registry.Tool, _ platform.Platform, _ *github.Client, st *state.State) (current, latest string, err error) {
	current = st.ToolVersion(tool.Name)
	return current, "git-managed", nil
}

func (s *ShellPluginInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	dest := expandHome(tool.CloneDest, p.HomeDir)

	if _, err := os.Stat(dest); err == nil {
		cmd := exec.Command("git", "-C", dest, "pull", "--ff-only")
		if err := utils.RunCmd(cmd); err != nil {
			os.RemoveAll(dest)
		} else {
			st.SetToolVersion(tool.Name, "git-managed")
			return s.runPostClone(tool, p, st)
		}
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	cmd := exec.Command("git", "clone", "--depth=1", tool.CloneURL, dest)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("git clone failed: %w", err)}
	}

	st.SetToolVersion(tool.Name, "git-managed")
	return s.runPostClone(tool, p, st)
}

func (s *ShellPluginInstaller) runPostClone(tool *registry.Tool, p platform.Platform, st *state.State) Result {
	switch tool.PostClone {
	case "spaceship":
		return s.postCloneSpaceship(tool, p)
	case "nvchad":
		return s.postCloneNvChad(tool, p)
	case "tpm":
		// TPM install_plugins runs after tmux.conf is deployed (in init flow)
		return Result{Tool: tool.Name, Version: "git-managed"}
	default:
		return Result{Tool: tool.Name, Version: "git-managed"}
	}
}

func (s *ShellPluginInstaller) postCloneSpaceship(tool *registry.Tool, p platform.Platform) Result {
	dest := expandHome(tool.CloneDest, p.HomeDir)
	linkSrc := filepath.Join(dest, "spaceship.zsh-theme")
	linkDst := filepath.Join(p.HomeDir, ".oh-my-zsh", "custom", "themes", "spaceship.zsh-theme")

	os.Remove(linkDst)
	if err := os.Symlink(linkSrc, linkDst); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("symlink failed: %w", err)}
	}

	return Result{Tool: tool.Name, Version: "git-managed"}
}

func (s *ShellPluginInstaller) postCloneNvChad(tool *registry.Tool, p platform.Platform) Result {
	dest := expandHome(tool.CloneDest, p.HomeDir)
	chadrcPath := filepath.Join(dest, "lua", "chadrc.lua")
	data, err := os.ReadFile(chadrcPath)
	if err == nil {
		patched := strings.Replace(string(data), `theme = "`, `theme = "catppuccin", -- `, 1)
		if string(data) == patched {
			patched = strings.Replace(string(data), "theme =", `theme = "catppuccin", transparency = true, --`, 1)
		}
		if err := os.WriteFile(chadrcPath, []byte(patched), 0644); err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("failed to patch chadrc.lua: %w", err)}
		}
	}

	return Result{Tool: tool.Name, Version: "git-managed"}
}

func expandHome(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

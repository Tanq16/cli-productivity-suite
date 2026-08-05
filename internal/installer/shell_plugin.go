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

func (s *ShellPluginInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	dest := expandHome(tool.CloneDest, p.HomeDir)

	if _, err := os.Stat(dest); err == nil {
		cmd := exec.Command("git", "-C", dest, "pull", "--ff-only")
		if err := utils.RunCmd(cmd); err != nil {
			return Result{Tool: tool.Name, Err: fmt.Errorf("git pull failed (local changes?): %w — remove %s manually to reclone", err, dest)}
		}
		st.SetToolVersion(tool.Name, "git-managed")
		return Result{Tool: tool.Name, Version: "git-managed"}
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return Result{Tool: tool.Name, Err: err}
	}

	cmd := exec.Command("git", "clone", "--depth=1", tool.CloneURL, dest)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("git clone failed: %w", err)}
	}

	st.SetToolVersion(tool.Name, "git-managed")
	return Result{Tool: tool.Name, Version: "git-managed"}
}

func expandHome(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

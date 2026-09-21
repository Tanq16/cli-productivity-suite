package installer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type PythonToolInstaller struct{}

func (t *PythonToolInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	if tool.PyTool == "" {
		return Result{Tool: tool.Name, Err: errors.New("no uv tool defined")}
	}

	env := p.CustomScriptEnv()
	cmd := envCommand(env, "uv", "tool", "install", "--force", tool.PyTool)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("uv tool install %s failed: %w", tool.PyTool, err)}
	}

	version := installedPyToolVersion(env, tool.PyTool)
	st.SetToolVersion(tool.Name, version)
	return Result{Tool: tool.Name, Version: version}
}

func installedPyToolVersion(env []string, pkg string) string {
	out, err := envCommand(env, "uv", "tool", "list").Output()
	if err != nil {
		return "uv-managed"
	}
	for line := range strings.SplitSeq(string(out), "\n") {
		name, version, found := strings.Cut(strings.TrimSpace(line), " v")
		if found && name == pkg {
			return version
		}
	}
	return "uv-managed"
}

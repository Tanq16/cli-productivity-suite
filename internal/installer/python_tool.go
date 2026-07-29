package installer

import (
	"errors"
	"fmt"

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

	cmd := envCommand(p.CustomScriptEnv(), "uv", "tool", "install", "--force", tool.PyTool)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("uv tool install %s failed: %w", tool.PyTool, err)}
	}

	st.SetToolVersion(tool.Name, "uv-managed")
	return Result{Tool: tool.Name, Version: "uv-managed"}
}

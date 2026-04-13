package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

type CustomScriptInstaller struct{}

func (c *CustomScriptInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	if tool.InstallCmd == "" {
		return Result{Tool: tool.Name, Err: fmt.Errorf("no install command defined")}
	}

	cmd := exec.Command("bash", "-c", tool.InstallCmd)
	cmd.Env = os.Environ()
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("custom install failed: %w", err)}
	}

	st.SetToolVersion(tool.Name, "custom-managed")
	return Result{Tool: tool.Name, Version: "custom-managed"}
}

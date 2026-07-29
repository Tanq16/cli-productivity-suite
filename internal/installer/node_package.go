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

type NodePackageInstaller struct{}

func (n *NodePackageInstaller) Install(tool *registry.Tool, p platform.Platform, _ *github.Client, st *state.State) Result {
	if tool.NodePkg == "" {
		return Result{Tool: tool.Name, Err: errors.New("no npm package defined")}
	}

	cmd := envCommand(p.CustomScriptEnv(), "npm", "install", "-g", tool.NodePkg)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("npm install %s failed: %w", tool.NodePkg, err)}
	}

	st.SetToolVersion(tool.Name, "npm-managed")
	return Result{Tool: tool.Name, Version: "npm-managed"}
}

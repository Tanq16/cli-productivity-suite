package installer

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"strings"

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

	args := append([]string{"install", "-g"}, strings.Fields(tool.NodePkg)...)
	cmd := envCommand(p.CustomScriptEnv(), "npm", args...)
	if err := utils.RunCmd(cmd); err != nil {
		return Result{Tool: tool.Name, Err: fmt.Errorf("npm install %s failed: %w", tool.NodePkg, err)}
	}

	version := installedNodeVersion(p.CustomScriptEnv(), primaryNodePkg(tool.NodePkg))
	st.SetToolVersion(tool.Name, version)
	return Result{Tool: tool.Name, Version: version}
}

func primaryNodePkg(spec string) string {
	name, _, _ := strings.Cut(spec, " ")
	if base, _, found := strings.CutLast(name, "@"); found && base != "" {
		return base
	}
	return name
}

func installedNodeVersion(env []string, pkg string) string {
	out, err := envCommand(env, "npm", "ls", "-g", pkg, "--json", "--depth=0").Output()
	if err != nil {
		return "npm-managed"
	}
	var listed struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal(out, &listed); err != nil {
		return "npm-managed"
	}
	if dep, ok := listed.Dependencies[pkg]; ok && dep.Version != "" {
		return dep.Version
	}
	return "npm-managed"
}

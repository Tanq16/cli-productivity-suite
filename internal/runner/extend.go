package runner

import (
	"fmt"
	"os"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func ExtendList() {
	packs := registry.AllExtensionPacks()
	for _, pack := range packs {
		tools := filterExtPackForPlatform(pack)
		utils.PrintInfo(fmt.Sprintf("%s — %s (%d tools)", pack.Name, pack.Description, len(tools)))
		if len(tools) > 0 {
			names := ""
			for i, t := range tools {
				if i > 0 {
					names += ", "
				}
				names += t.Name
			}
			utils.PrintGeneric("    " + names)
		}
	}
}

func Extend(packName string, ghToken string) {
	pack := registry.ExtensionPackByName(packName)
	if pack == nil {
		available := ""
		for _, p := range registry.AllExtensionPacks() {
			if available != "" {
				available += ", "
			}
			available += p.Name
		}
		utils.PrintFatal(fmt.Sprintf("unknown extension pack: %s (available: %s)", packName, available), nil)
	}

	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	gh := github.NewClient(ghToken)

	if err := os.MkdirAll(p.ShellExtDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellExtDir()), err)
	}

	tools := filterExtPackForPlatform(*pack)
	if len(tools) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: no tools for this platform", pack.Name))
		return
	}

	for _, t := range tools {
		if t.IsPrivate && ghToken == "" {
			utils.PrintFatal(fmt.Sprintf("tool %s is private — provide --gh-token", t.Name), nil)
		}
	}

	var hadErrors bool
	phaseName := fmt.Sprintf("Extension: %s", pack.Name)

	if len(tools) == 1 {
		tool := tools[0]
		inst := installer.Dispatch(tool.Kind)
		if inst == nil {
			utils.PrintFatal(fmt.Sprintf("no installer for kind: %s", tool.Kind), nil)
		}
		utils.PrintRunning("installing " + tool.Name)
		st.Remove(tool.Name)
		result := inst.Install(&tool, p, gh, st)
		utils.ClearLines(1)
		if result.Err != nil {
			utils.PrintFatal(fmt.Sprintf("%s: install failed", tool.Name), result.Err)
		}
		utils.PrintSuccess(fmt.Sprintf("%s: installed %s", tool.Name, result.Version))
	} else if len(tools) > 1 {
		hadErrors = runPhase(phaseName, tools, p, gh, st)
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	if hadErrors {
		utils.PrintWarn(fmt.Sprintf("extension pack %s finished with errors", pack.Name), nil)
	} else {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s complete!", pack.Name))
	}
}

func ExtendCheck(packName string, ghToken string) {
	pack := registry.ExtensionPackByName(packName)
	if pack == nil {
		utils.PrintFatal(fmt.Sprintf("unknown extension pack: %s", packName), nil)
	}

	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	gh := github.NewClient(ghToken)
	tools := filterExtPackForPlatform(*pack)

	var checkable []registry.Tool
	for _, t := range tools {
		if t.IsPrivate && ghToken == "" {
			continue
		}
		checkable = append(checkable, t)
	}

	if len(checkable) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: no tools to check", pack.Name))
		return
	}

	utils.PrintRunning(fmt.Sprintf("Checking extension pack: %s", pack.Name))
	results := checkTools(checkable, p, gh, st)
	utils.ClearLines(1)

	if len(results) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: everything up to date", pack.Name))
		return
	}

	utils.PrintInfo(fmt.Sprintf("Extension pack %s check complete", pack.Name))
	for _, r := range results {
		switch r.Status {
		case "update":
			utils.PrintIndentedWarn(fmt.Sprintf("%s: update available (%s → %s)", r.Tool.Name, r.Current, r.Latest), nil)
		case "error":
			utils.PrintIndentedError(fmt.Sprintf("%s: check failed", r.Tool.Name), r.Err)
		}
	}
}

func filterExtPackForPlatform(pack registry.ExtensionPack) []registry.Tool {
	p, err := platform.Detect()
	if err != nil {
		return pack.Tools
	}
	return filterPlatformTools(pack.Tools, p)
}

func ExtensionPackNames() []string {
	packs := registry.AllExtensionPacks()
	names := make([]string, len(packs))
	for i, p := range packs {
		names[i] = p.Name
	}
	return names
}

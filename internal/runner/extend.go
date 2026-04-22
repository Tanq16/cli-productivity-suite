package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tanq16/cli-productivity-suite/internal/configs"
	"github.com/tanq16/cli-productivity-suite/internal/custom"
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var customOnce sync.Once

func ensureCustomPacks() {
	customOnce.Do(func() {
		p, err := platform.Detect()
		if err != nil {
			return
		}
		extDir := filepath.Join(p.ConfigDir(), "extensions")
		packs, warnings := custom.LoadDir(extDir, registry.BuiltinPackNames())
		for _, w := range warnings {
			utils.PrintWarn(w, nil)
		}
		registry.LoadCustomPacks(packs)
	})
}

func ExtendList() {
	ensureCustomPacks()
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

func Extend(packName string, toolFilter []string, ghToken string) {
	ensureCustomPacks()
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

	allTools := filterExtPackForPlatform(*pack)
	if len(allTools) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: no tools for this platform", pack.Name))
		return
	}

	if len(toolFilter) > 0 {
		nameSet := make(map[string]bool, len(toolFilter))
		for _, n := range toolFilter {
			nameSet[n] = true
		}
		var filtered []registry.Tool
		for _, t := range allTools {
			if nameSet[t.Name] {
				filtered = append(filtered, t)
				delete(nameSet, t.Name)
			}
		}
		if len(nameSet) > 0 {
			var unknown []string
			for name := range nameSet {
				unknown = append(unknown, name)
			}
			utils.PrintFatal(fmt.Sprintf("tools not found in extension pack %s: %s", packName, strings.Join(unknown, ", ")), nil)
		}
		allTools = filtered
	}

	var tools []registry.Tool
	for _, t := range allTools {
		if t.IsPrivate && ghToken == "" {
			utils.PrintWarn(fmt.Sprintf("skipping %s (private repo, no --gh-token)", t.Name), nil)
			continue
		}
		tools = append(tools, t)
	}

	if len(tools) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: all tools require --gh-token, nothing to install", pack.Name))
		return
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
		if result.Skipped {
			utils.PrintSuccess(fmt.Sprintf("%s: already at %s", tool.Name, result.Version))
		} else if result.WasUpdated {
			utils.PrintSuccess(fmt.Sprintf("%s: updated to %s", tool.Name, result.Version))
		} else {
			utils.PrintSuccess(fmt.Sprintf("%s: installed %s", tool.Name, result.Version))
		}
	} else if len(tools) > 1 {
		hadErrors = runPhase(phaseName, tools, p, gh, st)
	}

	deployPackFragment(packName, p)

	var hasBinaries bool
	for _, t := range tools {
		switch t.Kind {
		case registry.GitHubRelease, registry.DirectDownload, registry.LanguageRuntime:
			hasBinaries = true
		}
	}
	if hasBinaries {
		var compErrors []jobResult
		var compLines int
		utils.PrintRunning("(Running) Regenerating completions")
		generateCompletions(p, &compErrors, &compLines)
		utils.ClearLines(compLines + 1)
		if len(compErrors) > 0 {
			utils.PrintError("Regenerating completions: partially completed with errors", nil)
			for _, e := range compErrors {
				utils.PrintIndentedError(e.name, e.err)
			}
		} else {
			utils.PrintInfo("Regenerating completions")
		}
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

func ExtendRemove(packName string, toolFilter []string) {
	ensureCustomPacks()

	if registry.BuiltinPackNames()[packName] {
		utils.PrintFatal(fmt.Sprintf("--remove is only supported for custom extension packs; %q is built-in", packName), nil)
	}

	pack := registry.ExtensionPackByName(packName)
	if pack == nil {
		utils.PrintFatal(fmt.Sprintf("unknown custom extension pack: %s", packName), nil)
	}

	p, err := platform.Detect()
	if err != nil {
		utils.PrintFatal("platform detection failed", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintFatal("failed to load state", err)
	}

	allTools := filterExtPackForPlatform(*pack)
	if len(allTools) == 0 {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: no tools for this platform", pack.Name))
		return
	}

	if len(toolFilter) > 0 {
		nameSet := make(map[string]bool, len(toolFilter))
		for _, n := range toolFilter {
			nameSet[n] = true
		}
		var filtered []registry.Tool
		for _, t := range allTools {
			if nameSet[t.Name] {
				filtered = append(filtered, t)
				delete(nameSet, t.Name)
			}
		}
		if len(nameSet) > 0 {
			var unknown []string
			for name := range nameSet {
				unknown = append(unknown, name)
			}
			utils.PrintFatal(fmt.Sprintf("tools not found in extension pack %s: %s", packName, strings.Join(unknown, ", ")), nil)
		}
		allTools = filtered
	}

	wholePackRemoval := len(toolFilter) == 0
	var hadErrors bool

	for _, t := range allTools {
		if t.RemoveCmd == "" {
			utils.PrintError(fmt.Sprintf("%s: no remove command defined in YAML", t.Name), nil)
			hadErrors = true
			continue
		}
		utils.PrintRunning("removing " + t.Name)
		cmd := exec.Command("bash", "-c", t.RemoveCmd)
		cmd.Env = os.Environ()
		err := utils.RunCmd(cmd)
		utils.ClearLines(1)
		if err != nil {
			utils.PrintError(fmt.Sprintf("%s: remove failed", t.Name), err)
			hadErrors = true
			continue
		}
		st.Remove(t.Name)
		utils.PrintSuccess(fmt.Sprintf("%s: removed", t.Name))
	}

	if wholePackRemoval {
		fragPath := filepath.Join(p.ShellDir(), "rc", "custom", packName+".zsh")
		if err := os.Remove(fragPath); err != nil && !os.IsNotExist(err) {
			utils.PrintError(fmt.Sprintf("failed to remove fragment %s", fragPath), err)
			hadErrors = true
		}
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
	}

	if hadErrors {
		utils.PrintWarn(fmt.Sprintf("extension pack %s: removal finished with errors", pack.Name), nil)
	} else {
		utils.PrintSuccess(fmt.Sprintf("extension pack %s: removal complete", pack.Name))
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
	ensureCustomPacks()
	packs := registry.AllExtensionPacks()
	names := make([]string, len(packs))
	for i, p := range packs {
		names[i] = p.Name
	}
	return names
}

func ExtensionPackToolNames(packName string) []string {
	ensureCustomPacks()
	pack := registry.ExtensionPackByName(packName)
	if pack == nil {
		return nil
	}
	tools := filterExtPackForPlatform(*pack)
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}

func deployPackFragment(packName string, p platform.Platform) {
	switch packName {
	case "runtimes":
		deployFragment(p, "10-runtimes.zsh", configs.RcRuntimes())
	case "cloud":
		deployFragment(p, "20-cloud.zsh", configs.RcCloud())
	case "security":
		content := []byte("export NUCLEI_TEMPLATES_DIR=\"$HOME/shell/nuclei-templates\"\n")
		deployFragment(p, "30-security.zsh", content)
	default:
		if pf := custom.GetPackFile(packName); pf != nil {
			content := custom.RenderFragment(*pf)
			if content != nil {
				deployFragment(p, filepath.Join("custom", packName+".zsh"), content)
			}
		}
	}
}

func deployFragment(p platform.Platform, filename string, content []byte) {
	dest := filepath.Join(p.ShellDir(), "rc", filename)
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		utils.PrintError(fmt.Sprintf("failed to create %s", filepath.Dir(dest)), err)
		return
	}
	if err := os.WriteFile(dest, content, 0644); err != nil {
		utils.PrintError(fmt.Sprintf("failed to write %s", dest), err)
	}
}

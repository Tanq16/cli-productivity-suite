package runner

import (
	"fmt"
	"os"
	"strings"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func PackageList() {
	p, _ := platformAndState()

	for _, g := range registry.Groups() {
		tools := filterPlatformTools(registry.ByGroup(g.Name), p)
		if len(tools) == 0 {
			continue
		}
		utils.PrintInfo(fmt.Sprintf("%s - %s (%d packages)", g.Name, g.Description, len(tools)))

		if g.Name != registry.GroupRuntime {
			utils.PrintIndentedList(toolNames(tools))
			continue
		}
		for _, s := range registry.Suites() {
			suiteTools := filterPlatformTools(registry.BySuite(s.Name), p)
			if len(suiteTools) == 0 {
				continue
			}
			utils.PrintIndentedGeneric(fmt.Sprintf("%s: %s", s.Name, strings.Join(toolNames(suiteTools), ", ")))
		}
	}
}

func PackageInstall(target, ghToken string) {
	p, st := platformAndState()

	resolved, ok := registry.Resolve(target)
	if !ok {
		utils.PrintFatal(fmt.Sprintf("unknown package, group or suite: %s (run `cps package list`)", target), nil)
	}

	tools := filterPlatformTools(resolved, p)
	if len(tools) == 0 {
		utils.PrintSuccess(fmt.Sprintf("%s: nothing to install on this platform", target))
		return
	}
	tools = filterPlatformTools(registry.Order(tools, st.Installed), p)

	if err := os.MkdirAll(p.ShellExtDir(), 0755); err != nil {
		utils.PrintFatal(fmt.Sprintf("failed to create %s", p.ShellExtDir()), err)
	}

	var skipped []string
	var installable []registry.Tool
	for _, t := range tools {
		if t.IsPrivate && ghToken == "" {
			skipped = append(skipped, t.Name)
			continue
		}
		installable = append(installable, t)
	}
	if len(skipped) > 0 {
		utils.PrintWarn(fmt.Sprintf("skipping %s (private repo, no --gh-token)", strings.Join(skipped, ", ")), nil)
	}
	if len(installable) == 0 {
		utils.PrintSuccess(fmt.Sprintf("%s: every package needs --gh-token, nothing to install", target))
		return
	}

	gh := github.NewClient(ghToken)
	hadErrors := runPhase("Installing "+target, installable, p, gh, st)

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
		hadErrors = true
	}

	if hasBinaries(installable) && runPostInstall("Regenerating completions", p, false) {
		hadErrors = true
	}
	if deployRCFile("Regenerating the rc file", p, st) {
		hadErrors = true
	}

	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
		hadErrors = true
	}

	if hadErrors {
		utils.PrintFatal(fmt.Sprintf("%s finished with errors", target), nil)
	}
	utils.PrintSuccess(fmt.Sprintf("%s complete!", target))
}

func toolNames(tools []registry.Tool) []string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}

func hasBinaries(tools []registry.Tool) bool {
	for _, t := range tools {
		switch t.Kind {
		case registry.GitHubRelease, registry.DirectDownload, registry.LanguageRuntime:
			return true
		}
	}
	return false
}

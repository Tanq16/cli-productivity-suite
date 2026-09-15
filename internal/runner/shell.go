package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/rcgen"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func Shell() {
	p, st := platformAndState()
	gh := github.NewClient("")

	if platform.BrewPath() == "" {
		utils.PrintWarn("Homebrew is not installed; run `cps prereq` first or the shell env will be incomplete", nil)
	}

	utils.PrintRunning("Phase 1: Shell directories")
	for _, dir := range []string{
		p.ShellDir(),
		filepath.Join(p.ShellDir(), "rc"),
		filepath.Join(p.ShellDir(), "rc", "custom"),
		filepath.Join(p.ShellDir(), "env"),
		filepath.Join(p.ShellDir(), "plugins"),
		filepath.Join(p.ShellDir(), "custom-bin"),
		p.ShellExtDir(),
		p.ShellAppsDir(),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			utils.ClearLines(1)
			utils.PrintFatal(fmt.Sprintf("failed to create %s", dir), err)
		}
	}
	utils.ClearLines(1)
	utils.PrintInfo("Phase 1: Shell directories")

	tools := filterPlatformTools(registry.ShellTools(), p)
	var hadErrors bool

	for _, phase := range []struct {
		name string
		kind registry.ToolKind
	}{
		{"Phase 2: Shell binaries", registry.GitHubRelease},
		{"Phase 3: Applications", registry.AppBundle},
		{"Phase 4: Shell plugins", registry.RepoSnapshot},
		{"Phase 5: Config files", registry.ConfigFile},
	} {
		if runPhase(phase.name, filterKind(tools, phase.kind), p, gh, st) {
			hadErrors = true
		}
		if err := st.Save(); err != nil {
			utils.PrintError("failed to save state", err)
			hadErrors = true
		}
	}

	if runPostInstall("Phase 6: Shell environment", p, true) {
		hadErrors = true
	}

	if deployRCFile("Phase 7: Shell rc file", p, st) {
		hadErrors = true
	}

	st.LastInit = time.Now()
	if err := st.Save(); err != nil {
		utils.PrintError("failed to save state", err)
		hadErrors = true
	}

	if hadErrors {
		utils.PrintFatal("shell setup finished with errors", nil)
	}
	utils.PrintSuccess("shell setup complete!")
}

func deployRCFile(phaseName string, p platform.Platform, st *state.State) bool {
	utils.PrintRunning(phaseName)
	err := rcgen.Generate(p, st)
	utils.ClearLines(1)
	if err != nil {
		utils.PrintError(phaseName+": failed", err)
		return true
	}
	utils.PrintInfo(phaseName)
	return false
}

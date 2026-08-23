package runner

import (
	"fmt"

	"github.com/tanq16/cli-productivity-suite/internal/migration"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

// MigrationGate runs before every command. It blocks once per upgrade so a user cannot act on a new
// cps against an environment the previous version laid out, and records the version once acknowledged.
func MigrationGate(appVersion string) {
	p, err := platform.Detect()
	if err != nil {
		return
	}
	st, err := state.Load(p.StatePath())
	if err != nil {
		return
	}
	if st.CPSVersion() == appVersion {
		return
	}

	// A first run has nothing laid out by an earlier version, so it records the version and never prompts.
	pending := migration.Pending(st.CPSVersion(), appVersion)
	if len(pending) == 0 || st.IsFresh() {
		recordVersion(st, appVersion)
		return
	}

	utils.PrintWarn(fmt.Sprintf("cps upgraded to %s and this needs manual migration first", appVersion), nil)
	for _, m := range pending {
		utils.PrintInfo("leaving " + m.From + " behind:")
		for _, note := range m.Notes {
			utils.PrintGeneric("    " + note)
		}
	}

	choice, err := utils.PromptSelect("Have you completed the steps above?", []string{"No, exit so I can run them", "Yes, continue"})
	if err != nil {
		utils.PrintFatal("could not read your answer", err)
	}
	if choice != 1 {
		utils.PrintFatal("run the steps above, then re-run this command", nil)
	}

	recordVersion(st, appVersion)
	utils.PrintSuccess("migration acknowledged, recorded " + appVersion)
}

func recordVersion(st *state.State, appVersion string) {
	st.SetCPSVersion(appVersion)
	if err := st.Save(); err != nil {
		utils.PrintWarn("could not record the cps version in state", err)
	}
}

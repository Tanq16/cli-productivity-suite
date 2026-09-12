package runner

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/status"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func Status(check, asJSON bool, ghToken string) {
	p, st := platformAndState()
	entries := status.Collect(p, st)

	if check {
		if !asJSON {
			utils.PrintRunning("checking upstream versions")
		}
		status.Check(entries, p, github.NewClient(ghToken))
		if !asJSON {
			utils.ClearLines(1)
		}
	}

	if asJSON {
		data, err := json.Marshal(entries, jsontext.WithIndent("  "))
		if err != nil {
			utils.PrintFatal("failed to render status as JSON", err)
		}
		utils.PrintGeneric(string(data))
		return
	}

	var installed, outdated, failed int
	for _, group := range status.GroupOrder() {
		inGroup := status.ByGroup(entries, group)
		if len(inGroup) == 0 {
			continue
		}
		utils.PrintInfo(group)
		for _, e := range inGroup {
			switch {
			case !e.Installed:
				utils.PrintGeneric("    " + e.Name + ": not installed")
			case e.CheckError != "":
				failed++
				installed++
				utils.PrintIndentedWarn(fmt.Sprintf("%s: %s (check failed)", e.Name, e.Version), nil)
			case e.Outdated:
				installed++
				outdated++
				utils.PrintIndentedWarn(fmt.Sprintf("%s: %s → %s", e.Name, e.Version, e.Latest), nil)
			default:
				installed++
				utils.PrintIndentedSuccess(fmt.Sprintf("%s: %s", e.Name, e.Version))
			}
		}
	}

	summary := fmt.Sprintf("%d of %d packages installed", installed, len(entries))
	if check {
		summary += fmt.Sprintf(", %d outdated", outdated)
		if failed > 0 {
			summary += fmt.Sprintf(", %d could not be checked", failed)
		}
	}
	utils.PrintSuccess(summary)
}

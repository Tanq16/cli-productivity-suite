package runner

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/internal/status"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func Status(check, asJSON bool, ghToken string) {
	p, err := platform.Detect()
	if err != nil {
		statusFatal(asJSON, "platform detection failed", err)
	}
	st, err := state.Load(p.StatePath())
	if err != nil {
		statusFatal(asJSON, "failed to load state", err)
	}

	entries := status.Collect(p, st)
	if check {
		runCheck(entries, p, github.NewClient(ghToken), asJSON)
	}

	if asJSON {
		data, err := json.Marshal(entries, jsontext.WithIndent("  "))
		if err != nil {
			statusFatal(asJSON, "failed to render status as JSON", err)
		}
		utils.PrintGeneric(string(data))
		return
	}

	var installed, outdated, failed, skipped int
	for _, group := range status.GroupOrder() {
		inGroup := status.ByGroup(entries, group)
		if len(inGroup) == 0 {
			continue
		}
		utils.PrintInfo(group)
		for _, e := range inGroup {
			switch {
			case !e.Installed:
				utils.PrintIndentedGeneric(e.Name + ": not installed")
			case e.CheckError != "":
				failed++
				installed++
				utils.PrintIndentedWarn(fmt.Sprintf("%s: %s (check failed)", e.Name, e.Version), errors.New(e.CheckError))
			case check && e.SkipReason != "":
				skipped++
				installed++
				utils.PrintIndentedGeneric(fmt.Sprintf("%s: %s (not checked: %s)", e.Name, e.Version, e.SkipReason))
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
		if skipped > 0 {
			summary += fmt.Sprintf(", %d not checked", skipped)
		}
		if failed > 0 {
			summary += fmt.Sprintf(", %d could not be checked", failed)
		}
	}
	utils.PrintSuccess(summary)
}

func runCheck(entries []status.Entry, p platform.Platform, gh *github.Client, asJSON bool) {
	if asJSON {
		for i := range entries {
			status.CheckEntry(&entries[i], p, gh)
		}
		return
	}

	var total int64
	for _, e := range entries {
		if e.Installed {
			total++
		}
	}

	m := utils.NewMeter("Checking", "upstream versions", total, packagesUnit)
	for i := range entries {
		if !entries[i].Installed {
			continue
		}
		m.Item(entries[i].Name)
		status.CheckEntry(&entries[i], p, gh)
		m.Add(1)
	}
	m.Done()
}

func statusFatal(asJSON bool, msg string, err error) {
	if !asJSON {
		utils.PrintFatal(msg, err)
	}

	payload := map[string]string{"error": msg}
	if err != nil {
		payload["detail"] = err.Error()
	}
	data, marshalErr := json.Marshal(payload, jsontext.WithIndent("  "))
	if marshalErr != nil {
		utils.PrintFatal(msg, err)
	}
	utils.PrintGeneric(string(data))
	os.Exit(1)
}

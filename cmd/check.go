package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var checkFlags struct {
	public  bool
	private bool
	system  bool
	all     bool
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	RunE:  runCheck,
}

func init() {
	checkCmd.Flags().BoolVar(&checkFlags.public, "public", false, "Check public tools only")
	checkCmd.Flags().BoolVar(&checkFlags.private, "private", false, "Check private tools only")
	checkCmd.Flags().BoolVar(&checkFlags.system, "system", false, "Check system packages only")
	checkCmd.Flags().BoolVar(&checkFlags.all, "all", false, "Check all tools")
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintError("platform detection failed", err)
		return fmt.Errorf("platform detection failed: %w", err)
	}

	st, err := state.Load(p.StatePath())
	if err != nil {
		utils.PrintError("failed to load state", err)
		return fmt.Errorf("failed to load state: %w", err)
	}

	gh := github.NewClient(ghToken)
	reg := registry.New()

	// Default to --all if no filter specified
	if !checkFlags.public && !checkFlags.private && !checkFlags.system {
		checkFlags.all = true
	}

	var tools []registry.Tool
	if checkFlags.all {
		tools = reg.ForPlatform(p.OS.String())
	} else {
		if checkFlags.public {
			tools = append(tools, reg.ByCategory(registry.Public)...)
		}
		if checkFlags.private {
			tools = append(tools, reg.ByCategory(registry.Private)...)
		}
		if checkFlags.system {
			tools = append(tools, reg.ByCategory(registry.System)...)
		}
	}

	// Only check tools that have state entries (i.e., are installed)
	headers := []string{"Tool", "Current", "Latest", "Status"}
	var rows [][]string

	for _, tool := range tools {
		current := st.ToolVersion(tool.Name)
		if current == "" {
			continue
		}

		t := tool
		inst := installer.Dispatch(t.Kind)
		if inst == nil {
			continue
		}

		cur, lat, err := inst.Check(&t, p, gh, st)
		if err != nil {
			rows = append(rows, []string{t.Name, cur, "error", err.Error()})
			continue
		}

		status := "up-to-date"
		if cur != lat && lat != "" && lat != "system-managed" && lat != "check-manually" && lat != "embedded" && lat != "git-managed" {
			status = "update available"
		}
		rows = append(rows, []string{t.Name, cur, lat, status})
	}

	if len(rows) == 0 {
		utils.PrintWarn("no installed tools found in state", nil)
		return nil
	}

	utils.PrintTable(headers, rows)
	return nil
}

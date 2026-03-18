package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/orchestrator"
)

var checkFlags struct {
	Public  bool
	Private bool
	System  bool
	All     bool
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := orchestrator.CheckConfig{
			Public:  checkFlags.Public,
			Private: checkFlags.Private,
			System:  checkFlags.System,
			All:     checkFlags.All,
		}
		orchestrator.RunCheck(ghToken, AppVersion, cfg)
	},
}

func init() {
	checkCmd.Flags().BoolVarP(&checkFlags.Public, "public", "p", false, "Check public tools only")
	checkCmd.Flags().BoolVarP(&checkFlags.Private, "private", "P", false, "Check private tools only")
	checkCmd.Flags().BoolVarP(&checkFlags.System, "system", "s", false, "Check system packages only")
	checkCmd.Flags().BoolVarP(&checkFlags.All, "all", "a", false, "Check all tools")
}

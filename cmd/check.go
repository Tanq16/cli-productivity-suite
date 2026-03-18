package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/orchestrator"
)

var checkFlags orchestrator.CheckFlags

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Run: func(cmd *cobra.Command, args []string) {
		orchestrator.RunCheck(ghToken, AppVersion, checkFlags)
	},
}

func init() {
	checkCmd.Flags().BoolVarP(&checkFlags.Public, "public", "p", false, "Check public tools only")
	checkCmd.Flags().BoolVarP(&checkFlags.Private, "private", "P", false, "Check private tools only")
	checkCmd.Flags().BoolVarP(&checkFlags.System, "system", "s", false, "Check system packages only")
	checkCmd.Flags().BoolVarP(&checkFlags.All, "all", "a", false, "Check all tools")
}

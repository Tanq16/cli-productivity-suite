package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var checkFlags struct {
	SkipPrivate bool
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for available updates",
	Run: func(cmd *cobra.Command, args []string) {
		runner.Check(ghToken, AppVersion, checkFlags.SkipPrivate)
	},
}

func init() {
	checkCmd.Flags().BoolVar(&checkFlags.SkipPrivate, "skip-private", false, "Skip checking private tools")
}

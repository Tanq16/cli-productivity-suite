package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var prereqFlags struct {
	excludeBrew bool
}

var prereqCmd = &cobra.Command{
	Use:   "prereq",
	Short: "Install the system prerequisites cps needs (uses sudo)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Prereq(prereqFlags.excludeBrew)
	},
}

func init() {
	prereqCmd.Flags().BoolVar(&prereqFlags.excludeBrew, "exclude-brew", false, "Skip the Homebrew install (Linux only; ignored on macOS)")
}

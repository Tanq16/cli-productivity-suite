package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update cps itself to the latest release",
	Run: func(cmd *cobra.Command, args []string) {
		runner.SelfUpdate(AppVersion)
	},
}

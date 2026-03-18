package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/orchestrator"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all CPS-managed files and directories",
	Run: func(cmd *cobra.Command, args []string) {
		orchestrator.RunClean()
	},
}

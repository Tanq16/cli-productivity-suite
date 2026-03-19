package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize full dev environment setup",
	Run: func(cmd *cobra.Command, args []string) {
		runner.Init(ghToken)
	},
}

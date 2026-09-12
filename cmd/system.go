package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var systemCmd = &cobra.Command{
	Use:       "system <verb>",
	Short:     "Run a system-level operation (uses sudo)",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: runner.SystemVerbs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.System(args[0])
	},
}

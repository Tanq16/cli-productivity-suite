package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var customCmd = &cobra.Command{
	Use:                "custom [args...]",
	Short:              "Run ~/.config/cps/custom.sh with the cps PATH",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Custom(args)
	},
}

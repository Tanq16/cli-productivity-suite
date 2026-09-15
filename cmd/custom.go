package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var customCmd = &cobra.Command{
	Use:   "custom [args...]",
	Short: "Run ~/.config/cps/custom.sh with the cps PATH",
	Long:  "Run ~/.config/cps/custom.sh with the cps PATH.\n\nArguments pass through verbatim. Put flag-shaped arguments after a -- separator,\nso that `cps custom --debug` stays the cps flag and `cps custom -- --debug` reaches\nthe script.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Custom(args)
	},
}

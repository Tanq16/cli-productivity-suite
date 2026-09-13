package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var themeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Switch the terminal colour theme (run without a name to list)",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			runner.ThemeList()
			return
		}
		runner.Theme(args[0])
	},
}

func init() {
	themeCmd.ValidArgs = runner.ThemeNames()
}

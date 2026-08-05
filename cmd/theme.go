package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var themeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Switch the terminal colour theme (run without a name to list)",
	Args:  cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return runner.ThemeNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			runner.ThemeList()
			return
		}
		runner.Theme(args[0])
	},
}

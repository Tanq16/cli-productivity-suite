package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var extendCheckFlag bool

var extendCmd = &cobra.Command{
	Use:   "extend <pack-name>",
	Short: "Install extension tool packs (e.g., cloud-sec, app-sec)",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return runner.ExtensionPackNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if extendCheckFlag {
			runner.ExtendCheck(args[0], ghToken)
		} else {
			runner.Extend(args[0], ghToken)
		}
	},
}

var extendListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available extension packs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.ExtendList()
	},
}

func init() {
	extendCmd.Flags().BoolVar(&extendCheckFlag, "check", false, "Check for updates instead of installing")
	extendCmd.AddCommand(extendListCmd)
}

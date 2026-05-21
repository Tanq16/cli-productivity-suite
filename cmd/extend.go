package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var extendFlags struct {
	remove bool
}

var extendCmd = &cobra.Command{
	Use:   "extend <pack-name> [tools...]",
	Short: "Install extension tool packs (e.g., security, cloudsec, runtimes)",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			suggestions := append([]string{"list"}, runner.ExtensionPackNames()...)
			return suggestions, cobra.ShellCompDirectiveNoFileComp
		}
		// For subsequent args, suggest tool names from the selected pack
		return runner.ExtensionPackToolNames(args[0]), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "list" {
			runner.ExtendList()
			return
		}
		if extendFlags.remove {
			runner.ExtendRemove(args[0], args[1:])
			return
		}
		runner.Extend(args[0], args[1:], ghToken)
	},
}

func init() {
	extendCmd.Flags().BoolVar(&extendFlags.remove, "remove", false, "Remove tool(s) from a custom extension pack (custom packs only)")
}

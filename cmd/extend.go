package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var extendFlags struct {
	check bool
}

var extendCmd = &cobra.Command{
	Use:   "extend <pack-name>",
	Short: "Install extension tool packs (e.g., security, cloudsec, appsec)",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		suggestions := append([]string{"list"}, runner.ExtensionPackNames()...)
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "list" {
			runner.ExtendList()
			return
		}
		if extendFlags.check {
			runner.ExtendCheck(args[0], ghToken)
		} else {
			runner.Extend(args[0], ghToken)
		}
	},
}

func init() {
	extendCmd.Flags().BoolVar(&extendFlags.check, "check", false, "Check for updates instead of installing")
}

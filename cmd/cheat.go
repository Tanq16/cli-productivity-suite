package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/cheatsheet"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var cheatCmd = &cobra.Command{
	Use:   "cheat <topic>",
	Short: "Print cheat sheets for common tools",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return cheatsheet.AllNames(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cheatsheet.Print(args[0]); err != nil {
			utils.PrintFatal(err.Error(), nil)
		}
	},
}

var cheatListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available cheat sheets",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, s := range cheatsheet.List() {
			utils.PrintInfo(fmt.Sprintf("%s — %s", s.Name, s.Description))
		}
	},
}

func init() {
	cheatCmd.AddCommand(cheatListCmd)
}

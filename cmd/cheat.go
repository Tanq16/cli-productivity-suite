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
		suggestions := append([]string{"list"}, cheatsheet.AllNames()...)
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "list" {
			for _, s := range cheatsheet.List() {
				utils.PrintInfo(fmt.Sprintf("%s — %s", s.Name, s.Description))
			}
			return
		}
		if err := cheatsheet.Print(args[0]); err != nil {
			utils.PrintFatal("failed to print cheat sheet", err)
		}
	},
}

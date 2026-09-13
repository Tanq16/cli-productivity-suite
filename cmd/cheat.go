package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/cheatsheet"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var cheatCmd = &cobra.Command{
	Use:       "cheat <topic>",
	Short:     "Print cheat sheets for common tools",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: append([]string{"list"}, cheatsheet.AllNames()...),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "list" {
			for _, s := range cheatsheet.List() {
				utils.PrintInfo(fmt.Sprintf("%s - %s", s.Name, s.Description))
			}
			return
		}
		if err := cheatsheet.Print(args[0]); err != nil {
			utils.PrintFatal("failed to print cheat sheet", err)
		}
	},
}

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/orchestrator"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
)

var installCmd = &cobra.Command{
	Use:               "install <tool>",
	Short:             "Install a single tool by name",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeToolNames,
	Run: func(cmd *cobra.Command, args []string) {
		orchestrator.RunInstall(args[0], ghToken)
	},
}

func completeToolNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	reg := registry.New()
	return reg.Names(), cobra.ShellCompDirectiveNoFileComp
}

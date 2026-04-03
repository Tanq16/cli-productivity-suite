package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var installCmd = &cobra.Command{
	Use:               "install <tools|categories>...",
	Short:             "Install tools by name or category",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completeInstallArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Install(args, ghToken)
	},
}

func completeInstallArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	suggestions := runner.CategoryAliasNames()
	suggestions = append(suggestions, registry.New().Names()...)
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

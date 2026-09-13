package packageCmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var installFlags struct {
	ghToken string
}

var PackageCmd = &cobra.Command{
	Use:   "package",
	Short: "List and install packages, groups and suites",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List every group, suite and package",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.PackageList()
	},
}

var installCmd = &cobra.Command{
	Use:       "install <package|group|suite>",
	Short:     "Install a package, a whole group, or a runtime suite",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: registry.InstallTargets(),
	Run: func(cmd *cobra.Command, args []string) {
		runner.PackageInstall(args[0], github.ResolveToken(installFlags.ghToken))
	},
}

func init() {
	PackageCmd.AddCommand(installCmd)
	PackageCmd.AddCommand(listCmd)

	installCmd.Flags().StringVar(&installFlags.ghToken, "gh-token", github.TokenFromEnv(), "GitHub PAT for private repos")
	installCmd.Flags().Lookup("gh-token").DefValue = "$GITHUB_TOKEN / $GH_TOKEN"
}

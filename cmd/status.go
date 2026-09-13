package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var statusFlags struct {
	check   bool
	asJSON  bool
	ghToken string
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show what is installed and, with --check, what is out of date",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Status(statusFlags.check, statusFlags.asJSON, github.ResolveToken(statusFlags.ghToken))
	},
}

func init() {
	statusCmd.Flags().BoolVar(&statusFlags.check, "check", false, "Compare every installed package against its upstream version")
	statusCmd.Flags().BoolVar(&statusFlags.asJSON, "json", false, "Emit the status as JSON")
	statusCmd.Flags().StringVar(&statusFlags.ghToken, "gh-token", github.TokenFromEnv(), "GitHub PAT, which raises the API rate limit and reaches private repos")
	statusCmd.Flags().Lookup("gh-token").DefValue = "$GITHUB_TOKEN / $GH_TOKEN"
}

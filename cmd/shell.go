package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Install the shell environment: binaries, Neovim, plugins and configs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runner.Shell()
	},
}

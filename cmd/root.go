package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-suite",
	Short: "CLI Productivity Suite installer and manager",
	Long: `A command line tool to install, update, fix and delete the CLI Productivity Suite.
Complete documentation is available at https://github.com/tanq16/cli-productivity-suite`,
}

func Execute() error {
	return rootCmd.Execute()
}

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-productivity-suite",
	Short: "CLI Productivity Suite manager that installs a funky shell and general tools for MacOS and Linux",
}

func Execute() error {
	return rootCmd.Execute()
}

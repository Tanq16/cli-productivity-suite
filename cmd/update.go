package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var updateFlags struct {
	IncludeConf bool
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed tools to latest versions",
	Run: func(cmd *cobra.Command, args []string) {
		runner.Update(ghToken, updateFlags.IncludeConf)
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateFlags.IncludeConf, "include-conf", "c", false, "Also overwrite deployed config files")
}

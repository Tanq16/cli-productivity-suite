package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/orchestrator"
)

var updateFlags orchestrator.UpdateFlags

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed tools to latest versions",
	Run: func(cmd *cobra.Command, args []string) {
		orchestrator.RunUpdate(cmd.Context(), ghToken, updateFlags)
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateFlags.Public, "public", "p", false, "Update public tools only")
	updateCmd.Flags().BoolVarP(&updateFlags.Private, "private", "P", false, "Update private tools only")
	updateCmd.Flags().BoolVarP(&updateFlags.System, "system", "s", false, "Update system packages only")
	updateCmd.Flags().BoolVarP(&updateFlags.All, "all", "a", false, "Update all tools")
	updateCmd.Flags().BoolVarP(&updateFlags.IncludeConf, "include-conf", "c", false, "Also overwrite deployed config files")
}

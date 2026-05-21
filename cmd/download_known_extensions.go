package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/runner"
)

var downloadKnownExtensionsCmd = &cobra.Command{
	Use:   "download-known-extensions",
	Short: "Download reference custom-extension YAMLs from the CPS repo",
	Long: `Fetches the reference custom-extension packs maintained in the CPS repo
(ai-tools, additional-cloud-tools, database, praetorian) and writes them to
~/.config/cps/extensions/, where they become available via cps extend list.

Overwrites existing files of the same name — if you've customized one of the
reference packs locally, rename it before re-running this.`,
	Run: func(cmd *cobra.Command, args []string) {
		runner.DownloadKnownExtensions()
	},
}

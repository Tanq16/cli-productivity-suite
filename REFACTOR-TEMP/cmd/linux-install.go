package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install CLI stuff for linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting installation")
		if err := installer.Install(); err != nil {
			log.Fatal().Err(err).Msg("Installation failed")
		}
		log.Info().Msg("Installation completed successfully")
	},
}

func init() {
	linuxCmd.AddCommand(installCmd)
}

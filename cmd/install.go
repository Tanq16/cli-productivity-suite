package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/internal/installer"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the CLI productivity suite",
	Long:  `Install all components of the CLI productivity suite including tools, configurations, and plugins`,
	RunE:  runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	log.Info().Msg("Starting installation")
	if err := installer.Install(); err != nil {
		log.Error().Err(err).Msg("Installation failed")
		return err
	}
	log.Info().Msg("Installation completed successfully")
	return nil
}

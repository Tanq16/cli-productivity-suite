package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/internal"
)

var rootCmd = &cobra.Command{
	Use:   "cli-suite",
	Short: "CLI Productivity Suite manager that installs a funky shell and general tools for MacOS and Linux",
}

var linuxCmd = &cobra.Command{
	Use:   "linux",
	Short: "CLI Productivity Suite - Linux operations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Linux operations")
	},
}

var linuxInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install CLI stuff for linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting installation")
		if err := internal.InstallLinux(); err != nil {
			log.Fatal().Err(err).Msg("Installation failed")
		}
		log.Info().Msg("Installation completed successfully")
	},
}

var linuxCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove CLI stuff for linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting installation")
		if err := internal.CleanLinux(); err != nil {
			log.Fatal().Err(err).Msg("Installation failed")
		}
		log.Info().Msg("Installation completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(linuxCmd)
	linuxCmd.AddCommand(linuxInstallCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

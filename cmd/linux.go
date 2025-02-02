package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var linuxCmd = &cobra.Command{
	Use:   "linux",
	Short: "CLI Productivity Suite - Linux operations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Linux operations")
	},
}

func init() {
	rootCmd.AddCommand(linuxCmd)
}

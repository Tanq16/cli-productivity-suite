package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var AppVersion = "dev-build"

var ghToken string
var debugFlag, forAIFlag bool

var rootCmd = &cobra.Command{
	Use:               "cps",
	Short:             "CLI Productivity Suite — manage your dev environment",
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = AppVersion
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&forAIFlag, "for-ai", false, "AI-friendly output (markdown tables, no color)")
	rootCmd.MarkFlagsMutuallyExclusive("debug", "for-ai")

	defaultToken := os.Getenv("CPS_GITHUB_PAT")
	rootCmd.PersistentFlags().StringVar(&ghToken, "gh-token", defaultToken, "GitHub PAT for private repos (env: CPS_GITHUB_PAT)")

	cobra.OnInitialize(setupLogs)
}

func setupLogs() {
	if debugFlag {
		utils.GlobalDebugFlag = true
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Str("package", "cmd").Msg("debug logging enabled")
	} else if forAIFlag {
		utils.GlobalForAIFlag = true
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

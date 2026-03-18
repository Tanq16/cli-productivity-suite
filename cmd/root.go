package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
	Version:           AppVersion,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&forAIFlag, "for-ai", false, "AI-friendly output (markdown tables, no color)")
	rootCmd.MarkFlagsMutuallyExclusive("debug", "for-ai")

	rootCmd.PersistentFlags().StringVar(&ghToken, "gh-token", "", "GitHub PAT for private repos")

	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)

	cobra.OnInitialize(setupLogs, resolveGHToken)
}

func resolveGHToken() {
	if ghToken != "" {
		return
	}
	out, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		return
	}
	ghToken = strings.TrimSpace(string(out))
}

func setupLogs() {
	if debugFlag {
		utils.GlobalDebugFlag = true
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
			NoColor:    false,
		}
		log.Logger = zerolog.New(output).With().Timestamp().Logger()
		log.Debug().Str("package", "cmd").Msg("debug logging enabled")
	} else if forAIFlag {
		utils.GlobalForAIFlag = true
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}

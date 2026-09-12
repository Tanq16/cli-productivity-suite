package cmd

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	packageCmd "github.com/tanq16/cli-productivity-suite/cmd/package-cmd"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var AppVersion = "dev-build"

var debugFlag bool

var rootCmd = &cobra.Command{
	Use:               "cps",
	Short:             "Manage your CLI dev environment",
	Version:           AppVersion,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")

	rootCmd.AddCommand(cheatCmd)
	rootCmd.AddCommand(customCmd)
	rootCmd.AddCommand(packageCmd.PackageCmd)
	rootCmd.AddCommand(prereqCmd)
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(themeCmd)

	cobra.OnInitialize(setupLogs)
}

func setupLogs() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var out io.Writer = os.Stdout
	if utils.StdoutIsTerminal {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	}
	log.Logger = zerolog.New(out).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		utils.GlobalDebugFlag = true
	}
}

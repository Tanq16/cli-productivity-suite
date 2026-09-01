package cmd

import (
	"cmp"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/internal/runner"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var AppVersion = "dev-build"

var ghToken string
var debugFlag, migrationAcknowledged bool

var rootCmd = &cobra.Command{
	Use:               "cps",
	Short:             "Manage your CLI dev environment",
	Version:           AppVersion,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		runner.MigrationGate(AppVersion, migrationAcknowledged)
		resolveGHToken()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&migrationAcknowledged, "migration-acknowledged", false, "Confirm the manual migration steps an upgrade printed have been run")

	defaultGHToken := cmp.Or(os.Getenv("GITHUB_TOKEN"), os.Getenv("GH_TOKEN"))
	for _, c := range []*cobra.Command{initCmd, extendCmd} {
		c.Flags().StringVar(&ghToken, "gh-token", defaultGHToken, "GitHub PAT for private repos (or GITHUB_TOKEN / GH_TOKEN env)")
	}

	rootCmd.AddCommand(cheatCmd)
	rootCmd.AddCommand(extendCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(themeCmd)

	cobra.OnInitialize(setupLogs)
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

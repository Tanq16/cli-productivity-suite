package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all CPS-managed files and directories",
	RunE:  runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	p, err := platform.Detect()
	if err != nil {
		utils.PrintError("platform detection failed", err)
		return fmt.Errorf("platform detection failed: %w", err)
	}

	dirs := []string{
		filepath.Join(p.HomeDir, "shell"),
		filepath.Join(p.HomeDir, ".tmux"),
		filepath.Join(p.HomeDir, ".config", "nvim"),
		filepath.Join(p.HomeDir, ".nvm"),
		filepath.Join(p.HomeDir, "nuclei-templates"),
		filepath.Join(p.HomeDir, "google-cloud-sdk"),
		filepath.Join(p.HomeDir, ".config", "cps"),
	}

	utils.PrintWarn("this will remove the following directories:", nil)
	for _, d := range dirs {
		utils.PrintGeneric("  " + d)
	}
	utils.PrintGeneric("\nare you sure? (yes/no): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))

	if answer != "yes" {
		utils.PrintInfo("clean aborted")
		return nil
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			continue
		}
		if err := os.RemoveAll(d); err != nil {
			utils.PrintError(fmt.Sprintf("failed to remove %s", d), err)
		} else {
			utils.PrintSuccess(fmt.Sprintf("removed %s", d))
		}
	}

	utils.PrintSuccess("clean complete")
	return nil
}

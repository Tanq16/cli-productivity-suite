package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tanq16/cli-productivity-suite/utils"
)

var privilegedCmd = &cobra.Command{
	Use:    "_privileged",
	Short:  "Run privileged operations (internal use only)",
	Hidden: true,
}

var privInstallPkgsCmd = &cobra.Command{
	Use:   "install-packages",
	Short: "Install system packages via apt",
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.PrintInfo("updating apt cache")
		update := exec.Command("apt-get", "update", "-y")
		if err := update.Run(); err != nil {
			return fmt.Errorf("apt update failed: %w", err)
		}

		pkgs := []string{
			"tmux", "git", "tree", "wget", "curl", "zsh",
			"openssl", "nmap", "ncat",
			"cmake", "gcc", "make", "ninja-build", "gettext", "unzip", "file",
			"neovim", "ffmpeg",
		}

		utils.PrintInfo("installing packages via apt")
		installArgs := append([]string{"install", "-y"}, pkgs...)
		install := exec.Command("apt-get", installArgs...)
		if err := install.Run(); err != nil {
			return fmt.Errorf("apt install failed: %w", err)
		}

		utils.PrintSuccess("system packages installed")
		return nil
	},
}

func init() {
	privilegedCmd.AddCommand(privInstallPkgsCmd)
	rootCmd.AddCommand(privilegedCmd)
}

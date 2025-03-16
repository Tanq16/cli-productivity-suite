package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tanq16/cli-productivity-suite/utils"
)

func InstallOhMyZsh() error {
	homeDir := os.Getenv("HOME")
	omzDir := filepath.Join(homeDir, ".oh-my-zsh")
	if _, err := os.Stat(omzDir); err == nil {
		log.Info().Msg("Oh My Zsh is already installed")
		return nil
	}
	installerURL := "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
	installerPath := filepath.Join(os.TempDir(), "install_omz.sh")
	if err := utils.DownloadFile(installerURL, installerPath); err != nil {
		return fmt.Errorf("failed to download Oh My Zsh installer: %w", err)
	}
	defer os.Remove(installerPath)
	if err := os.Chmod(installerPath, 0755); err != nil {
		return fmt.Errorf("failed to make installer executable: %w", err)
	}
	if err := utils.ExecuteInteractiveCommand(installerPath); err != nil {
		return fmt.Errorf("failed to run Oh My Zsh installer: %w", err)
	}
	if err := installOhMyZshPlugins(); err != nil {
		return fmt.Errorf("failed to install Oh My Zsh plugins: %w", err)
	}
	return nil
}

func installOhMyZshPlugins() error {
	homeDir := os.Getenv("HOME")
	customDir := filepath.Join(homeDir, ".oh-my-zsh/custom")

	spaceshipDir := filepath.Join(customDir, "themes/spaceship-prompt")
	if _, err := utils.ExecuteCommand("git", "clone", "--depth=1",
		"https://github.com/spaceship-prompt/spaceship-prompt.git",
		spaceshipDir); err != nil {
		return fmt.Errorf("failed to clone spaceship prompt: %w", err)
	}
	spaceshipTheme := filepath.Join(customDir, "themes/spaceship.zsh-theme")
	if err := os.Symlink(
		filepath.Join(spaceshipDir, "spaceship.zsh-theme"),
		spaceshipTheme,
	); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create spaceship theme symlink: %w", err)
	}

	plugins := []struct {
		name string
		url  string
	}{
		{
			name: "zsh-autosuggestions",
			url:  "https://github.com/zsh-users/zsh-autosuggestions",
		},
		{
			name: "zsh-syntax-highlighting",
			url:  "https://github.com/zsh-users/zsh-syntax-highlighting",
		},
	}
	for _, plugin := range plugins {
		pluginDir := filepath.Join(customDir, "plugins", plugin.name)
		if _, err := utils.ExecuteCommand("git", "clone", "--depth=1",
			plugin.url, pluginDir); err != nil {
			return fmt.Errorf("failed to clone plugin %s: %w", plugin.name, err)
		}
	}
	return nil
}

func InstallNeovim() error {
	homeDir := os.Getenv("HOME")
	configPaths := []string{
		filepath.Join(homeDir, ".vim"),
		filepath.Join(homeDir, ".vimrc"),
		filepath.Join(homeDir, ".config/nvim"),
		filepath.Join(homeDir, ".local/share/nvim"),
	}
	// TODO: Remove deletion and default to only installing
	for _, path := range configPaths {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			log.Warn().Str("path", path).Msg("Failed to remove existing configuration")
		}
	}
	nvimConfigDir := filepath.Join(homeDir, ".config/nvim")
	if err := os.MkdirAll(filepath.Dir(nvimConfigDir), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if _, err := utils.ExecuteCommand("git", "clone", "--depth=1",
		"https://github.com/NvChad/starter", nvimConfigDir); err != nil {
		return fmt.Errorf("failed to clone NvChad configuration: %w", err)
	}
	if _, err := utils.ExecuteCommand("nvim", "--headless", "-c", "quitall"); err != nil {
		log.Warn().Msg("Initial Neovim startup produced errors (this is often normal)")
	}
	return nil
}

func InstallTmux() error {
	homeDir := os.Getenv("HOME")
	tpmDir := filepath.Join(homeDir, ".tmux/plugins/tpm")
	if err := os.MkdirAll(filepath.Dir(tpmDir), 0755); err != nil {
		return fmt.Errorf("failed to create tmux plugin directory: %w", err)
	}
	if _, err := utils.ExecuteCommand("git", "clone", "--depth=1",
		"https://github.com/tmux-plugins/tpm", tpmDir); err != nil {
		return fmt.Errorf("failed to clone TPM: %w", err)
	}
	tmuxConfPath := filepath.Join(homeDir, ".tmux.conf")
	if err := utils.DownloadFile(
		"https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/LEGACY/tmuxconf",
		tmuxConfPath,
	); err != nil {
		return fmt.Errorf("failed to download tmux configuration: %w", err)
	}
	pluginInstallCmd := filepath.Join(tpmDir, "bin/install_plugins")
	if _, err := utils.ExecuteCommand(pluginInstallCmd); err != nil {
		log.Warn().Msg("Tmux plugin installation produced errors (this is often normal)")
	}
	return nil
}

func ConfigureShell() error {
	homeDir := os.Getenv("HOME")
	zshrcPath := filepath.Join(homeDir, ".zshrc")

	// Backup existing configuration
	if _, err := os.Stat(zshrcPath); err == nil {
		backupPath := zshrcPath + ".backup"
		if err := os.Rename(zshrcPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing zshrc: %w", err)
		}
		log.Info().Str("backup", backupPath).Msg("Created backup of existing zshrc")
	}

	// Create new configuration
	rcFile, err := os.Create(zshrcPath)
	if err != nil {
		return fmt.Errorf("failed to create new zshrc: %w", err)
	}
	defer rcFile.Close()
	baseConfig := []string{
		"export ZSH=\"$HOME/.oh-my-zsh\"",
		"ZSH_THEME=\"spaceship\"",
		"plugins=(git zsh-autosuggestions zsh-syntax-highlighting)",
		"source $ZSH/oh-my-zsh.sh",
	}
	for _, line := range baseConfig {
		if _, err := rcFile.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write base configuration: %w", err)
		}
	}
	if err := mergeCustomRC(rcFile); err != nil {
		return fmt.Errorf("failed to merge custom RC configuration: %w", err)
	}
	return nil
}

func mergeCustomRC(rcFile *os.File) error {
	// TODO: check if this is the correct way to check for the OS
	rcFileName := "linux.rcfile"
	if os.Getenv("GOOS") == "darwin" {
		rcFileName = "macos.rcfile"
	}
	tempRCPath := filepath.Join(os.TempDir(), rcFileName)
	if err := utils.DownloadFile(
		fmt.Sprintf("https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/LEGACY/%s", rcFileName),
		tempRCPath,
	); err != nil {
		return fmt.Errorf("failed to download RC file: %w", err)
	}
	defer os.Remove(tempRCPath)

	customRC, err := os.Open(tempRCPath)
	if err != nil {
		return fmt.Errorf("failed to open custom RC file: %w", err)
	}
	defer customRC.Close()
	scanner := bufio.NewScanner(customRC)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && line != "" {
			if _, err := rcFile.WriteString(line + "\n"); err != nil {
				return fmt.Errorf("failed to write custom RC line: %w", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read custom RC file: %w", err)
	}
	return nil
}

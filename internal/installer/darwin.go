package installer

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tanq16/cli-productivity-suite/internal/config"
	"github.com/tanq16/cli-productivity-suite/internal/utils"
)

type DarwinInstaller struct {
	BaseInstaller
}

func NewDarwinInstaller() *DarwinInstaller {
	installer := &DarwinInstaller{}
	installer.initializeSteps()
	return installer
}

func (i *DarwinInstaller) Install() error {
	for _, step := range i.Steps {
		log.Info().Str("step", step.Name).Msg("Starting installation step")
		if err := step.Execute(); err != nil {
			return fmt.Errorf("failed to execute step '%s': %w", step.Name, err)
		}
		log.Info().Str("step", step.Name).Msg("Completed installation step")
	}

	return nil
}

func (i *DarwinInstaller) initializeSteps() {
	i.Steps = []InstallStep{
		{
			Name:    "Install packages",
			Execute: i.InstallPackages,
		},
		{
			Name:    "Install Aerospace",
			Execute: i.installAerospace,
		},
		{
			Name:    "Install Oh My Zsh",
			Execute: i.InstallOhMyZsh,
		},
		{
			Name:    "Install Neovim",
			Execute: i.InstallNeovim,
		},
		{
			Name:    "Install tmux",
			Execute: i.InstallTmux,
		},
		{
			Name:    "Configure shell",
			Execute: i.ConfigureShell,
		},
	}
}

func (i *DarwinInstaller) InstallPackages() error {
	for _, pkg := range config.MacOSPackages {
		_, err := utils.ExecuteCommand("brew", "install", pkg)
		if err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}
	return nil
}

func (i *DarwinInstaller) installAerospace() error {
	_, err := utils.ExecuteCommand("brew", "install", "--cask", "nikitabobko/tap/aerospace")
	return err
}

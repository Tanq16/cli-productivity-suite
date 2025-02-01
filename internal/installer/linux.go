package installer

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tanq16/cli-productivity-suite/internal/config"
	"github.com/tanq16/cli-productivity-suite/internal/system"
	"github.com/tanq16/cli-productivity-suite/internal/utils"
)

type LinuxInstaller struct {
	BaseInstaller
}

func NewLinuxInstaller() *LinuxInstaller {
	installer := &LinuxInstaller{}
	installer.initializeSteps()
	return installer
}

func (i *LinuxInstaller) Install() error {
	if err := system.EnsureRoot(); err != nil {
		return fmt.Errorf("failed to get root privileges: %w", err)
	}

	for _, step := range i.Steps {
		log.Info().Str("step", step.Name).Msg("Starting installation step")
		if err := step.Execute(); err != nil {
			return fmt.Errorf("failed to execute step '%s': %w", step.Name, err)
		}
		log.Info().Str("step", step.Name).Msg("Completed installation step")
	}

	return nil
}

func (i *LinuxInstaller) initializeSteps() {
	i.Steps = []InstallStep{
		{
			Name:    "Update package lists",
			Execute: i.updatePackageLists,
		},
		{
			Name:    "Install packages",
			Execute: i.InstallPackages,
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

func (i *LinuxInstaller) updatePackageLists() error {
	_, err := utils.ExecuteCommand("apt", "update", "-y")
	return err
}

func (i *LinuxInstaller) InstallPackages() error {
	args := append([]string{"install", "-y"}, config.LinuxPackages...)
	_, err := utils.ExecuteCommand("apt", args...)
	return err
}

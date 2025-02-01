package installer

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tanq16/cli-productivity-suite/internal/system"
)

// Install performs the complete installation process
func Install() error {
	sysInfo := system.GetSystemInfo()
	log.Info().Str("os", sysInfo.OS).Str("arch", sysInfo.Arch).Msg("Detected system information")

	installer, err := getInstaller(sysInfo)
	if err != nil {
		return err
	}

	return installer.Install()
}

// getInstaller returns the appropriate installer for the system
func getInstaller(sysInfo system.SystemInfo) (PlatformInstaller, error) {
	switch sysInfo.OS {
	case "linux":
		return NewLinuxInstaller(), nil
	case "darwin":
		return NewDarwinInstaller(), nil
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", sysInfo.OS)
	}
}

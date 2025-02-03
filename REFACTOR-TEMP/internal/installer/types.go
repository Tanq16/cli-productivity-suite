package installer

// PlatformInstaller defines the interface for platform-specific installation
type PlatformInstaller interface {
	Install() error
	InstallPackages() error
	InstallOhMyZsh() error
	InstallNeovim() error
	InstallTmux() error
	ConfigureShell() error
}

// BaseInstaller contains common installation functionality
type BaseInstaller struct {
	Steps []InstallStep
}

type InstallStep struct {
	Name    string
	Execute func() error
}

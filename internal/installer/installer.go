package installer

import (
	"github.com/tanq16/cli-productivity-suite/internal/github"
	"github.com/tanq16/cli-productivity-suite/internal/platform"
	"github.com/tanq16/cli-productivity-suite/internal/registry"
	"github.com/tanq16/cli-productivity-suite/internal/state"
)

type Result struct {
	Tool       string
	Version    string
	WasUpdated bool
	Skipped    bool
	Err        error
}

type Installer interface {
	Install(tool *registry.Tool, p platform.Platform, gh *github.Client, st *state.State) Result
}

func Dispatch(kind registry.ToolKind) Installer {
	switch kind {
	case registry.GitHubRelease:
		return &GitHubReleaseInstaller{}
	case registry.DirectDownload:
		return &DirectDownloadInstaller{}
	case registry.SystemPackage:
		return &SystemPackageInstaller{}
	case registry.LanguageRuntime:
		return &LanguageRuntimeInstaller{}
	case registry.ConfigFile:
		return &ConfigDeployInstaller{}
	case registry.ShellPlugin:
		return &ShellPluginInstaller{}
	case registry.CustomScript:
		return &CustomScriptInstaller{}
	default:
		return nil
	}
}

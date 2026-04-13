package registry

type ToolKind int

const (
	GitHubRelease ToolKind = iota
	DirectDownload
	SystemPackage
	LanguageRuntime
	ConfigFile
	ShellPlugin
	CustomScript
)

func (k ToolKind) String() string {
	switch k {
	case GitHubRelease:
		return "github-release"
	case DirectDownload:
		return "direct-download"
	case SystemPackage:
		return "system-package"
	case LanguageRuntime:
		return "language-runtime"
	case ConfigFile:
		return "config-file"
	case ShellPlugin:
		return "shell-plugin"
	case CustomScript:
		return "custom-script"
	default:
		return "unknown"
	}
}

type ToolCategory int

const (
	Core ToolCategory = iota
	Private
	System
	Runtime
	Config
	Shell
	ExtSecurity
	ExtCloudSec
	ExtAppSec
	ExtMisc
	ExtPrivate
	ExtSystem
	ExtCloud
	ExtRuntimes
	ExtCustom
)

func (c ToolCategory) String() string {
	switch c {
	case Core:
		return "core"
	case Private:
		return "private"
	case System:
		return "system"
	case Runtime:
		return "runtime"
	case Config:
		return "config"
	case Shell:
		return "shell"
	case ExtSecurity:
		return "ext-security"
	case ExtCloudSec:
		return "ext-cloud-sec"
	case ExtAppSec:
		return "ext-app-sec"
	case ExtMisc:
		return "ext-misc"
	case ExtPrivate:
		return "ext-private"
	case ExtSystem:
		return "ext-system"
	case ExtCloud:
		return "ext-cloud"
	case ExtRuntimes:
		return "ext-runtimes"
	case ExtCustom:
		return "ext-custom"
	default:
		return "unknown"
	}
}

type AssetPattern struct {
	OSPatterns          map[string]string // "linux" -> "linux", "darwin" -> "apple" etc.
	ArchPatterns        map[string]string // "amd64" -> "x86_64", "arm64" -> "aarch64" etc.
	RequiredSubstrings  []string
	ExcludeSubstrings   []string
	ArchiveFormat       string            // "tar.gz", "tar.xz", "zip", "none" (raw binary)
	OSArchiveFormats    map[string]string // per-OS override, e.g. "linux" -> "tar.gz", "darwin" -> "zip"
	BinaryPathInArchive string            // glob pattern to find binary in extracted archive, e.g. "*/bat"
}

type Tool struct {
	Name        string
	BinaryName  string // name of the binary in ~/shell/executables/ (or ~/shell/extensions/ for extensions)
	Kind        ToolKind
	Category    ToolCategory
	Extension   bool   // true = install to ~/shell/extensions/ instead of ~/shell/executables/
	Repo        string // "owner/repo" for GitHub tools
	Asset       AssetPattern
	IsPrivate   bool
	BrewPkgs    []string // Homebrew packages (Linux + macOS)
	BrewCasks   []string // macOS brew cask packages
	Platforms   []string // "linux", "darwin", or both; empty means both
	Description string
	URL         string // for DirectDownload: URL template with {version}, {os}, {arch}
	StableURL   string // for DirectDownload: URL to fetch latest stable version string
	CloneURL    string // for ShellPlugin: full git clone URL
	CloneDest   string // for ShellPlugin: destination path (can use ~ for home)
	PostClone   string // for ShellPlugin: identifier for post-clone hook logic
	InstallCmd  string // for CustomScript: shell command run via bash -c
}

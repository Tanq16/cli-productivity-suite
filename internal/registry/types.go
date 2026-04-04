package registry

type ToolKind int

const (
	GitHubRelease ToolKind = iota
	DirectDownload
	SystemPackage
	CloudCLI
	LanguageRuntime
	ConfigFile
	ShellPlugin
)

func (k ToolKind) String() string {
	switch k {
	case GitHubRelease:
		return "github-release"
	case DirectDownload:
		return "direct-download"
	case SystemPackage:
		return "system-package"
	case CloudCLI:
		return "cloud-cli"
	case LanguageRuntime:
		return "language-runtime"
	case ConfigFile:
		return "config-file"
	case ShellPlugin:
		return "shell-plugin"
	default:
		return "unknown"
	}
}

type ToolCategory int

const (
	Public ToolCategory = iota
	Private
	System
	CloudCLICat
	Runtime
	Config
	Shell
	ExtCloudSec
	ExtAppSec
)

func (c ToolCategory) String() string {
	switch c {
	case Public:
		return "public"
	case Private:
		return "private"
	case System:
		return "system"
	case CloudCLICat:
		return "cloud-cli"
	case Runtime:
		return "runtime"
	case Config:
		return "config"
	case Shell:
		return "shell"
	case ExtCloudSec:
		return "ext-cloud-sec"
	case ExtAppSec:
		return "ext-app-sec"
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
	AptPkgs     []string // Linux apt packages
	BrewPkgs    []string // macOS brew packages
	BrewCasks   []string // macOS brew cask packages
	Platforms   []string // "linux", "darwin", or both; empty means both
	Description string
	URL         string // for DirectDownload: URL template with {version}, {os}, {arch}
	StableURL   string // for DirectDownload: URL to fetch latest stable version string
	CloneURL    string // for ShellPlugin: full git clone URL
	CloneDest   string // for ShellPlugin: destination path (can use ~ for home)
	PostClone   string // for ShellPlugin: identifier for post-clone hook logic
	ConfigSrc   string // for ConfigFile: which embedded config to use
	ConfigDest  string // for ConfigFile: destination path
}

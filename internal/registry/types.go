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
	NodePackage
	PythonTool
	AppBundle
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
	case NodePackage:
		return "node-package"
	case PythonTool:
		return "python-tool"
	case AppBundle:
		return "app-bundle"
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
	ExtEssentials
	ExtSecurity
	ExtCloudSec
	ExtAppSec
	ExtMisc
	ExtPrivate
	ExtSystem
	ExtCloud
	ExtRuntimes
	ExtAITools
	ExtHomelab
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
	case ExtEssentials:
		return "ext-essentials"
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
	case ExtAITools:
		return "ext-ai-tools"
	case ExtHomelab:
		return "ext-homelab"
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
	BinaryName  string // name of the installed binary in ~/shell/extensions/
	Kind        ToolKind
	Category    ToolCategory
	Extension   bool   // false = base tool installed by cps init; true = installed by cps extend
	Repo        string // "owner/repo" for GitHub tools
	Asset       AssetPattern
	IsPrivate   bool
	BrewPkgs    []string // Homebrew packages (Linux + macOS)
	BrewCasks   []string // macOS brew cask packages
	Platforms   []string // "linux", "darwin", or both; empty means both
	Description string
	URL         string // for DirectDownload and URL-sourced AppBundle: URL template with {version}, {os}, {arch}
	StableURL   string // for DirectDownload: URL to fetch latest stable version string
	CloneURL    string // for ShellPlugin: full git clone URL
	CloneDest   string // for ShellPlugin: destination path (can use ~ for home)
	InstallCmd  string // for CustomScript: shell command run via bash -c
	NodePkg     string // for NodePackage: npm package spec, e.g. "@openai/codex"
	PyTool      string // for PythonTool: uv tool name, e.g. "prowler"
	Requires    string // display-only prerequisite pack name; not enforced at install time

	VersionRegex string // applied to a StableURL body that is not a bare version string; first capture group wins
	PostInstall  string // for AppBundle: identifier for post-install hook logic
}

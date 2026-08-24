package registry

type ToolKind int

const (
	GitHubRelease ToolKind = iota
	DirectDownload
	SystemPackage
	LanguageRuntime
	ConfigFile
	RepoSnapshot
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
	case RepoSnapshot:
		return "repo-snapshot"
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

type AssetPattern struct {
	AssetNames          map[string]string // exact asset name keyed "<os>/<arch>"; a platform found here wins over every pattern below
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
	Dest        string // for RepoSnapshot: destination directory (can use ~ for home)
	InstallCmd  string // for CustomScript: shell command run via bash -c
	NodePkg     string // for NodePackage: npm package spec, e.g. "@openai/codex"
	PyTool      string // for PythonTool: uv tool name, e.g. "prowler"

	VersionRegex string // applied to a StableURL body that is not a bare version string; first capture group wins
	PostInstall  string // for AppBundle: identifier for post-install hook logic
}

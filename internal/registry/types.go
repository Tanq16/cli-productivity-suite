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
	AssetNames          map[string]string
	OSPatterns          map[string]string
	ArchPatterns        map[string]string
	RequiredSubstrings  []string
	ExcludeSubstrings   []string
	ArchiveFormat       string
	OSArchiveFormats    map[string]string
	BinaryPathInArchive string
}

type Tool struct {
	Name        string
	BinaryName  string
	Kind        ToolKind
	Extension   bool
	Repo        string
	Asset       AssetPattern
	IsPrivate   bool
	BrewPkgs    []string
	BrewCasks   []string
	Platforms   []string
	Description string
	URL         string
	StableURL   string
	Dest        string
	InstallCmd  string
	NodePkg     string
	PyTool      string

	VersionRegex string
	PostInstall  string
}

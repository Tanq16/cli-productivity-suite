package registry

// ExtensionPack defines a named collection of extension tools.
type ExtensionPack struct {
	Name        string
	Description string
	Category    ToolCategory
	Tools       []Tool
}

var extensionPacks = []ExtensionPack{
	{
		Name:        "security",
		Description: "General security tools",
		Category:    ExtSecurity,
		Tools: []Tool{
			{
				Name: "nuclei", BinaryName: "nuclei", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/nuclei", Description: "Vulnerability scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "nuclei",
				},
			},
			{
				Name: "naabu", BinaryName: "naabu", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/naabu", Description: "Port scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "naabu",
				},
			},
			{
				Name: "subfinder", BinaryName: "subfinder", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/subfinder", Description: "Subdomain discovery",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "subfinder",
				},
			},
			{
				Name: "proxify", BinaryName: "proxify", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/proxify", Description: "HTTP proxy",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "proxify",
				},
			},
			{
				Name: "trufflehog", BinaryName: "trufflehog", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "trufflesecurity/trufflehog", Description: "Secret scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "trufflehog",
				},
			},
		},
	},
	{
		Name:        "cloudsec",
		Description: "Cloud security and infrastructure tools",
		Category:    ExtCloudSec,
		Tools: []Tool{
			{
				Name: "kubelogin", BinaryName: "kubelogin", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "Azure/kubelogin", Description: "Azure Kubernetes login",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "bin/*/kubelogin",
				},
			},
			{
				Name: "grpcurl", BinaryName: "grpcurl", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "fullstorydev/grpcurl", Description: "curl for gRPC",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "osx"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "grpcurl",
				},
			},
			{
				Name: "terraform", BinaryName: "terraform", Kind: DirectDownload, Category: ExtCloudSec, Extension: true,
				Repo: "hashicorp/terraform", Description: "Infrastructure as code",
				URL:  "https://releases.hashicorp.com/terraform/{version_bare}/terraform_{version_bare}_{os}_{arch}.zip",
				Asset: AssetPattern{
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "terraform",
				},
			},
			{
				Name: "kubectl", BinaryName: "kubectl", Kind: DirectDownload, Category: ExtCloudSec, Extension: true,
				Description: "Kubernetes CLI",
				StableURL:   "https://dl.k8s.io/release/stable.txt",
				URL:         "https://dl.k8s.io/release/{version}/bin/{os}/{arch}/kubectl",
			},
		},
	},
	{
		Name:        "appsec",
		Description: "Application security tools",
		Category:    ExtAppSec,
		Tools: []Tool{
			{
				Name: "katana", BinaryName: "katana", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "projectdiscovery/katana", Description: "Web crawler",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "katana",
				},
			},
			{
				Name: "ffuf", BinaryName: "ffuf", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "ffuf/ffuf", Description: "Fast web fuzzer",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "ffuf",
				},
			},
		},
	},
	{
		Name:        "misc",
		Description: "Miscellaneous utility tools",
		Category:    ExtMisc,
		Tools:       []Tool{},
	},
	{
		Name:        "private",
		Description: "Private tools (requires --gh-token)",
		Category:    ExtPrivate,
		Tools: []Tool{
			{
				Name: "nits", BinaryName: "nits", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/nits", Description: "Nits tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "raikiri", BinaryName: "raikiri", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/raikiri", Description: "Raikiri tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "gcli", BinaryName: "gcli", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/gcli", Description: "Gcli tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "box", BinaryName: "box", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/box-cli", Description: "Box CLI tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "claudex", BinaryName: "claudex", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/claudex", Description: "Claudex tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "toon", BinaryName: "toon", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/toon", Description: "Private Toon tool", IsPrivate: true,
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "toon-*",
				},
			},
			{
				Name: "cybernest", BinaryName: "cybernest", Kind: GitHubRelease, Category: ExtPrivate, Extension: true,
				Repo: "Tanq16/cybernest", Description: "Private Cybernest tool", IsPrivate: true,
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
		},
	},
}

// AllExtensionPacks returns the list of available extension packs.
func AllExtensionPacks() []ExtensionPack {
	return extensionPacks
}

// ExtensionPackByName returns the pack with the given name, or nil.
func ExtensionPackByName(name string) *ExtensionPack {
	for i := range extensionPacks {
		if extensionPacks[i].Name == name {
			return &extensionPacks[i]
		}
	}
	return nil
}

// AllExtensionTools returns every tool across all extension packs.
func AllExtensionTools() []Tool {
	var all []Tool
	for _, pack := range extensionPacks {
		all = append(all, pack.Tools...)
	}
	return all
}

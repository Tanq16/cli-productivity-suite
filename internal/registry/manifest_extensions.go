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
			{
				Name: "httpx", BinaryName: "httpx", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/httpx", Description: "HTTP toolkit",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "httpx",
				},
			},
			{
				Name: "dnsx", BinaryName: "dnsx", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "projectdiscovery/dnsx", Description: "DNS toolkit",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "dnsx",
				},
			},
			{
				Name: "titus", BinaryName: "titus", Kind: GitHubRelease, Category: ExtSecurity, Extension: true,
				Repo: "praetorian-inc/titus", Description: "Security assessment tool",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{"windows", ".zip", ".jar", "browser-extension"},
					ArchiveFormat:     "none",
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
			{
				Name: "cloudfox", BinaryName: "cloudfox", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "BishopFox/cloudfox", Description: "Cloud security enumeration",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "cloudfox/cloudfox",
				},
			},
			{
				Name: "aurelian", BinaryName: "aurelian", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "praetorian-inc/aurelian", Description: "Cloud security tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", "SHA256SUMS"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "aurelian",
				},
			},
			{
				Name: "trivy", BinaryName: "trivy", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "aquasecurity/trivy", Description: "Vulnerability and misconfiguration scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "Linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "64bit", "arm64": "ARM64"},
					ExcludeSubstrings:   []string{"checksums", ".deb", ".rpm", "FreeBSD", "windows", "PPC64LE", "s390x", "bom.json", ".sigstore"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "trivy",
				},
			},
			{
				Name: "cloudlist", BinaryName: "cloudlist", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "projectdiscovery/cloudlist", Description: "Cloud asset discovery",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "windows"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "cloudlist",
				},
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
			{
				Name: "hadrian", BinaryName: "hadrian", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "praetorian-inc/hadrian", Description: "Application security scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "windows"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "hadrian",
				},
			},
			{
				Name: "dalfox", BinaryName: "dalfox", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "hahwul/dalfox", Description: "XSS scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "dalfox",
				},
			},
			{
				Name: "reaper", BinaryName: "reaper", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "ghostsecurity/reaper", Description: "API security testing",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "sigstore"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "reaper",
				},
			},
			{
				Name: "poltergeist", BinaryName: "poltergeist", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "ghostsecurity/poltergeist", Description: "API security tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"sigstore", "windows", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "poltergeist",
				},
			},
			{
				Name: "wraith", BinaryName: "wraith", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "ghostsecurity/wraith", Description: "API security tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "sigstore", "windows", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "wraith",
				},
			},
			{
				Name: "gau", BinaryName: "gau", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "lc/gau", Description: "URL fetcher",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "windows", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "gau",
				},
			},
		},
	},
	{
		Name:        "misc",
		Description: "Miscellaneous utility tools",
		Category:    ExtMisc,
		Tools: []Tool{
			{
				Name: "julius", BinaryName: "julius", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "praetorian-inc/julius", Description: "AI security testing tool",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{"checksums", ".exe"},
					ArchiveFormat:     "none",
				},
			},
			{
				Name: "trajan", BinaryName: "trajan", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "praetorian-inc/trajan", Description: "Security tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", ".zip", "windows"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "trajan",
				},
			},
			{
				Name: "gowitness", BinaryName: "gowitness", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "sensepost/gowitness", Description: "Web screenshot tool",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{"windows"},
					ArchiveFormat:     "none",
				},
			},
			{
				Name: "snitch", BinaryName: "snitch", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "karol-broda/snitch", Description: "Secret scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", ".deb", ".rpm", ".apk", "windows"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "snitch",
				},
			},
		},
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

package registry

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
				Name: "nuclei-templates", Kind: ShellPlugin, Category: ExtSecurity, Extension: true,
				Description: "Nuclei vulnerability templates",
				CloneURL:    "https://github.com/projectdiscovery/nuclei-templates.git",
				CloneDest:   "~/shell/nuclei-templates",
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
				URL: "https://releases.hashicorp.com/terraform/{version_bare}/terraform_{version_bare}_{os}_{arch}.zip",
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
				Name: "prowler", Kind: PythonTool, Category: ExtCloudSec, Extension: true,
				Description: "Cloud security posture scanner",
				PyTool:      "prowler", Requires: "runtimes",
			},
			{
				Name: "oci-cli", Kind: PythonTool, Category: ExtCloudSec, Extension: true,
				Description: "Oracle Cloud Infrastructure CLI",
				PyTool:      "oci-cli", Requires: "runtimes",
			},
			{
				Name: "tofu", BinaryName: "tofu", Kind: GitHubRelease, Category: ExtCloudSec, Extension: true,
				Repo: "opentofu/opentofu", Description: "OpenTofu infrastructure as code",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", "freebsd", "openbsd", "solaris", "sig", "pem", "SHA256SUMS", ".deb", ".rpm", ".apk", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "tofu",
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
				Name: "dalfox", BinaryName: "dalfox", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "hahwul/dalfox", Description: "XSS scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"windows", ".zip", ".sha256", "checksum", ".xml"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/dalfox",
				},
			},
			{
				Name: "gobuster", BinaryName: "gobuster", Kind: GitHubRelease, Category: ExtAppSec, Extension: true,
				Repo: "OJ/gobuster", Description: "Directory/DNS brute-forcer",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "Linux", "darwin": "Darwin"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "gobuster",
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
				Name: "age", BinaryName: "age", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "FiloSottile/age", Description: "File encryption tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", "freebsd"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "age/age",
				},
			},
			{
				Name: "sq", BinaryName: "sq", Kind: GitHubRelease, Category: ExtMisc, Extension: true,
				Repo: "neilotoole/sq", Description: "jq-like data wrangler for databases",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", ".deb", ".rpm", ".apk", ".zst", "checksums"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "sq",
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
	{
		Name:        "essentials",
		Description: "Core CLI binaries and starship prompt config",
		Category:    ExtEssentials,
		Tools: []Tool{
			{
				Name: "bat", BinaryName: "bat", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "sharkdp/bat", Description: "Cat clone with syntax highlighting",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"musl"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/bat",
				},
			},
			{
				Name: "fd", BinaryName: "fd", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "sharkdp/fd", Description: "Simple fast alternative to find",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"musl"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/fd",
				},
			},
			{
				Name: "ripgrep", BinaryName: "rg", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "BurntSushi/ripgrep", Description: "Fast recursive grep",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/rg",
				},
			},
			{
				Name: "lsd", BinaryName: "lsd", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "lsd-rs/lsd", Description: "Next gen ls command",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"musl"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/lsd",
				},
			},
			{
				Name: "jq", BinaryName: "jq", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "jqlang/jq", Description: "Command-line JSON processor",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "yq", BinaryName: "yq", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "mikefarah/yq", Description: "YAML processor",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{".tar.gz", ".zip"},
					ArchiveFormat:     "none",
				},
			},
			{
				Name: "fzf", BinaryName: "fzf", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "junegunn/fzf", Description: "Fuzzy finder",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "fzf",
				},
			},
			{
				Name: "gh", BinaryName: "gh", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "cli/cli", Description: "GitHub CLI",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					OSArchiveFormats:    map[string]string{"linux": "tar.gz", "darwin": "zip"},
					BinaryPathInArchive: "*/bin/gh",
				},
			},
			{
				Name: "gron", BinaryName: "gron", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "tomnomnom/gron", Description: "Make JSON greppable",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tgz",
					BinaryPathInArchive: "gron",
				},
			},
			{
				Name: "zoxide", BinaryName: "zoxide", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "ajeetdsouza/zoxide", Description: "Smarter cd command",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"android", ".deb"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "zoxide",
				},
			},
			{
				Name: "sd", BinaryName: "sd", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "chmln/sd", Description: "Find and replace CLI tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"windows", "gnueabi"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/sd",
				},
			},
			{
				Name: "starship", BinaryName: "starship", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "starship/starship", Description: "Minimal, fast, cross-shell prompt",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"gnu", "freebsd", "musleabihf", "i686"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "starship",
				},
			},
			{
				Name: "anbu", BinaryName: "anbu", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "Tanq16/anbu", Description: "Anbu tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "anbu*",
				},
			},
			{
				Name: "danzo", BinaryName: "danzo", Kind: GitHubRelease, Category: ExtEssentials, Extension: true,
				Repo: "Tanq16/danzo", Description: "Danzo tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "starship-config", Kind: ConfigFile, Category: ExtEssentials, Extension: true,
				Description: "Starship prompt configuration",
			},
		},
	},
	{
		Name:        "core",
		Description: "Dev tools, network utils, and media packages",
		Category:    ExtSystem,
		Tools: []Tool{
			{
				Name: "dev-tools", Kind: SystemPackage, Category: ExtSystem, Extension: true,
				Description: "Development build tools",
				Platforms:   []string{"linux"},
				BrewPkgs:    []string{"cmake", "gcc", "make", "ninja", "gettext"},
			},
			{
				Name: "network-tools", Kind: SystemPackage, Category: ExtSystem, Extension: true,
				Description: "Network utilities",
				BrewPkgs:    []string{"nmap", "openssl"},
			},
			{
				Name: "media-tools", Kind: SystemPackage, Category: ExtSystem, Extension: true,
				Description: "Media and monitoring tools",
				BrewPkgs:    []string{"ffmpeg"},
			},
			{
				Name: "aerospace", Kind: SystemPackage, Category: ExtSystem, Extension: true,
				Description: "macOS tiling window manager",
				Platforms:   []string{"darwin"},
				BrewCasks:   []string{"nikitabobko/tap/aerospace"},
			},
		},
	},
	{
		Name:        "cloud",
		Description: "Cloud CLIs (AWS, Azure, GCP)",
		Category:    ExtCloud,
		Tools: []Tool{
			{
				Name: "aws-cli", Kind: SystemPackage, Category: ExtCloud, Extension: true,
				Description: "AWS CLI v2",
				BrewPkgs:    []string{"awscli"},
			},
			{
				Name: "azure-cli", Kind: SystemPackage, Category: ExtCloud, Extension: true,
				Description: "Azure CLI",
				BrewPkgs:    []string{"azure-cli"},
			},
			{
				Name: "gcloud-cli", Kind: SystemPackage, Category: ExtCloud, Extension: true,
				Description: "Google Cloud CLI",
				BrewCasks:   []string{"gcloud-cli"},
			},
		},
	},
	{
		Name:        "runtimes",
		Description: "Language runtimes (Go, Rust, Python, Node)",
		Category:    ExtRuntimes,
		Tools: []Tool{
			{
				Name: "uv", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Python package manager (binary only)",
			},
			{
				Name: "fnm", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Node version manager (binary only)",
			},
			{
				Name: "bun", BinaryName: "bun", Kind: GitHubRelease, Category: ExtRuntimes, Extension: true,
				Repo: "oven-sh/bun", Description: "JavaScript runtime",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "x64", "arm64": "aarch64"},
					ExcludeSubstrings:   []string{"profile", "baseline", "musl"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "*/bun",
				},
			},
			{
				Name: "go-sdk", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Go programming language SDK",
			},
			{
				Name: "java-sdk", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Eclipse Temurin JDK (latest LTS)",
			},
			{
				Name: "python", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Python via uv + py-default venv",
			},
			{
				Name: "rust", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Rust toolchain via rustup",
			},
			{
				Name: "node", Kind: LanguageRuntime, Category: ExtRuntimes, Extension: true,
				Description: "Node.js LTS via fnm",
			},
		},
	},
	{
		Name:        "ai-tools",
		Description: "AI coding agents and LLM CLIs",
		Category:    ExtAITools,
		Tools: []Tool{
			{
				Name: "claude-code", Kind: NodePackage, Category: ExtAITools, Extension: true,
				Description: "Anthropic Claude Code CLI agent",
				NodePkg:     "@anthropic-ai/claude-code", Requires: "runtimes",
			},
			{
				Name: "antigravity", Kind: CustomScript, Category: ExtAITools, Extension: true,
				Description: "Google Antigravity CLI agent",
				InstallCmd: `set -eo pipefail
DEST_DIR="$HOME/shell/extensions"
mkdir -p "$DEST_DIR"
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux) OS=linux ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac
MANIFEST_URL="https://antigravity-cli-auto-updater-974169037036.us-central1.run.app/manifests/${OS}_${ARCH}.json"
URL=$(curl -fsSL "$MANIFEST_URL" | sed -n 's/.*"url"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
[ -n "$URL" ] || { echo "failed to parse antigravity manifest" >&2; exit 1; }
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
curl -fsSL "$URL" -o "$TMP/agy.tar.gz"
tar -xzf "$TMP/agy.tar.gz" -C "$TMP" antigravity
install -m 0755 "$TMP/antigravity" "$DEST_DIR/agy"
[ "$OS" = "darwin" ] && xattr -d com.apple.quarantine "$DEST_DIR/agy" 2>/dev/null || true
`,
			},
			{
				Name: "codex", Kind: NodePackage, Category: ExtAITools, Extension: true,
				Description: "OpenAI Codex CLI agent",
				NodePkg:     "@openai/codex", Requires: "runtimes",
			},
			{
				Name: "cursor-agent", Kind: CustomScript, Category: ExtAITools, Extension: true,
				Description: "Cursor headless coding agent",
				InstallCmd: `set -eo pipefail
DEST_DIR="$HOME/shell/extensions"
mkdir -p "$DEST_DIR"
case "$(uname -s)" in
  Darwin) OS=darwin ;;
  Linux) OS=linux ;;
  *) echo "unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64|amd64) ARCH=x64 ;;
  arm64|aarch64) ARCH=arm64 ;;
  *) echo "unsupported arch: $(uname -m)" >&2; exit 1 ;;
esac
VERSION=$(curl -fsSL https://cursor.com/install \
  | sed -n 's|.*downloads\.cursor\.com/lab/\([^/]*\)/.*|\1|p' | head -1)
[ -n "$VERSION" ] || { echo "failed to resolve cursor-agent version" >&2; exit 1; }
APP_DIR="$HOME/.local/share/cursor-agent"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
curl -fsSL "https://downloads.cursor.com/lab/${VERSION}/${OS}/${ARCH}/agent-cli-package.tar.gz" \
  | tar --strip-components=1 -xzf - -C "$TMP"
mkdir -p "$(dirname "$APP_DIR")"
rm -rf "$APP_DIR"
mv "$TMP" "$APP_DIR"
ln -sf "$APP_DIR/cursor-agent" "$DEST_DIR/cursor-agent"
`,
			},
		},
	},
	{
		Name:        "homelab",
		Description: "Self-hosted homelab services",
		Category:    ExtHomelab,
		Tools: []Tool{
			{
				Name: "caddy", BinaryName: "caddy", Kind: GitHubRelease, Category: ExtHomelab, Extension: true,
				Repo: "caddyserver/caddy", Description: "Web server / reverse proxy",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "mac"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"windows", "freebsd"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "caddy",
				},
			},
			{
				Name: "linksnapper", BinaryName: "linksnapper", Kind: GitHubRelease, Category: ExtHomelab, Extension: true,
				Repo: "Tanq16/linksnapper", Description: "LinkSnapper bookmark manager",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "kairo", BinaryName: "kairo", Kind: GitHubRelease, Category: ExtHomelab, Extension: true,
				Repo: "Tanq16/kairo", Description: "Kairo markdown note-taking app",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "raikiri", BinaryName: "raikiri", Kind: GitHubRelease, Category: ExtHomelab, Extension: true,
				Repo: "Tanq16/raikiri", Description: "Self-hosted media and music server",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "expenseowl", BinaryName: "expenseowl", Kind: GitHubRelease, Category: ExtHomelab, Extension: true,
				Repo: "Tanq16/expenseowl", Description: "Self-hosted expense tracker",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "rinnegan", Kind: AppBundle, Category: ExtHomelab, Extension: true,
				Repo: "Tanq16/rinnegan", Description: "Self-contained PTY web terminal",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "x64", "arm64": "arm64"},
					ArchiveFormat: "tar.gz",
				},
			},
			{
				Name: "code-server", Kind: AppBundle, Category: ExtHomelab, Extension: true,
				Repo: "coder/code-server", Description: "VS Code in the browser",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{".deb", ".rpm"},
					ArchiveFormat:     "tar.gz",
				},
			},
		},
	},
}

func AllExtensionPacks() []ExtensionPack {
	return extensionPacks
}

func ExtensionPackByName(name string) *ExtensionPack {
	for i := range extensionPacks {
		if extensionPacks[i].Name == name {
			return &extensionPacks[i]
		}
	}
	return nil
}

func AllExtensionTools() []Tool {
	var all []Tool
	for _, pack := range extensionPacks {
		all = append(all, pack.Tools...)
	}
	return all
}

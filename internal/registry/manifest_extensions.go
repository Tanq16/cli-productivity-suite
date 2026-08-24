package registry

type ExtensionPack struct {
	Name        string
	Description string
	Requires    string // display-only prerequisite pack name; not enforced at install time
	Tools       []Tool
}

var extensionPacks = []ExtensionPack{
	{
		Name:        "security",
		Description: "Security and application testing tools",
		Tools: []Tool{
			{
				Name: "nuclei", BinaryName: "nuclei", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/nuclei", Description: "Vulnerability scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "nuclei",
				},
			},
			{
				Name: "naabu", BinaryName: "naabu", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/naabu", Description: "Port scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "naabu",
				},
			},
			{
				Name: "subfinder", BinaryName: "subfinder", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/subfinder", Description: "Subdomain discovery",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "subfinder",
				},
			},
			{
				Name: "proxify", BinaryName: "proxify", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/proxify", Description: "HTTP proxy",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "proxify",
				},
			},
			{
				Name: "trufflehog", BinaryName: "trufflehog", Kind: GitHubRelease, Extension: true,
				Repo: "trufflesecurity/trufflehog", Description: "Secret scanner",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "trufflehog",
				},
			},
			{
				Name: "httpx", BinaryName: "httpx", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/httpx", Description: "HTTP toolkit",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "httpx",
				},
			},
			{
				Name: "dnsx", BinaryName: "dnsx", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/dnsx", Description: "DNS toolkit",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "dnsx",
				},
			},
			{
				Name: "nuclei-templates", Kind: RepoSnapshot, Extension: true,
				Description: "Nuclei vulnerability templates",
				Repo:        "projectdiscovery/nuclei-templates",
				Dest:        "~/shell/nuclei-templates",
			},
			{
				Name: "katana", BinaryName: "katana", Kind: GitHubRelease, Extension: true,
				Repo: "projectdiscovery/katana", Description: "Web crawler",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "katana",
				},
			},
			{
				Name: "ffuf", BinaryName: "ffuf", Kind: GitHubRelease, Extension: true,
				Repo: "ffuf/ffuf", Description: "Fast web fuzzer",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "ffuf",
				},
			},
			{
				Name: "dalfox", BinaryName: "dalfox", Kind: GitHubRelease, Extension: true,
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
				Name: "gobuster", BinaryName: "gobuster", Kind: GitHubRelease, Extension: true,
				Repo: "OJ/gobuster", Description: "Directory/DNS brute-forcer",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "Linux", "darwin": "Darwin"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "gobuster",
				},
			},
			{
				Name: "gau", BinaryName: "gau", Kind: GitHubRelease, Extension: true,
				Repo: "lc/gau", Description: "URL fetcher",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings:   []string{"checksums", "windows", ".zip"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "gau",
				},
			},
			{
				Name: "gowitness", BinaryName: "gowitness", Kind: GitHubRelease, Extension: true,
				Repo: "sensepost/gowitness", Description: "Web screenshot tool",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{"windows"},
					ArchiveFormat:     "none",
				},
			},
		},
	},
	{
		Name:        "cloudsec",
		Description: "Cloud security and infrastructure tools",
		Requires:    "runtimes",
		Tools: []Tool{
			{
				Name: "kubelogin", BinaryName: "kubelogin", Kind: GitHubRelease, Extension: true,
				Repo: "Azure/kubelogin", Description: "Azure Kubernetes login",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "bin/*/kubelogin",
				},
			},
			{
				Name: "grpcurl", BinaryName: "grpcurl", Kind: GitHubRelease, Extension: true,
				Repo: "fullstorydev/grpcurl", Description: "curl for gRPC",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "osx"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "grpcurl",
				},
			},
			{
				Name: "terraform", BinaryName: "terraform", Kind: DirectDownload, Extension: true,
				Repo: "hashicorp/terraform", Description: "Infrastructure as code",
				URL: "https://releases.hashicorp.com/terraform/{version_bare}/terraform_{version_bare}_{os}_{arch}.zip",
				Asset: AssetPattern{
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "terraform",
				},
			},
			{
				Name: "kubectl", BinaryName: "kubectl", Kind: DirectDownload, Extension: true,
				Description: "Kubernetes CLI",
				StableURL:   "https://dl.k8s.io/release/stable.txt",
				URL:         "https://dl.k8s.io/release/{version}/bin/{os}/{arch}/kubectl",
			},
			{
				Name: "trivy", BinaryName: "trivy", Kind: GitHubRelease, Extension: true,
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
				Name: "prowler", Kind: PythonTool, Extension: true,
				Description: "Cloud security posture scanner",
				PyTool:      "prowler",
			},
			{
				Name: "oci-cli", Kind: PythonTool, Extension: true,
				Description: "Oracle Cloud Infrastructure CLI",
				PyTool:      "oci-cli",
			},
			{
				Name: "tofu", BinaryName: "tofu", Kind: GitHubRelease, Extension: true,
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
		Name:        "private",
		Description: "Private tools (requires --gh-token)",
		Tools: []Tool{
			{
				Name: "nits", BinaryName: "nits", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/nits", Description: "Nits tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "gcli", BinaryName: "gcli", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/gcli", Description: "Gcli tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "box", BinaryName: "box", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/box-cli", Description: "Box CLI tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "claudex", BinaryName: "claudex", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/claudex", Description: "Claudex tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "sharingan", BinaryName: "sharingan", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/sharingan", Description: "EC2 workstation manager",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "toon", BinaryName: "toon", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/toon", Description: "Private Toon tool", IsPrivate: true,
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "toon-*",
				},
			},
			{
				Name: "cybernest", BinaryName: "cybernest", Kind: GitHubRelease, Extension: true,
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
		Tools: []Tool{
			{
				Name: "bat", BinaryName: "bat", Kind: GitHubRelease, Extension: true,
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
				Name: "fd", BinaryName: "fd", Kind: GitHubRelease, Extension: true,
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
				Name: "ripgrep", BinaryName: "rg", Kind: GitHubRelease, Extension: true,
				Repo: "BurntSushi/ripgrep", Description: "Fast recursive grep",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
					ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "*/rg",
				},
			},
			{
				Name: "lsd", BinaryName: "lsd", Kind: GitHubRelease, Extension: true,
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
				Name: "jq", BinaryName: "jq", Kind: GitHubRelease, Extension: true,
				Repo: "jqlang/jq", Description: "Command-line JSON processor",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "yq", BinaryName: "yq", Kind: GitHubRelease, Extension: true,
				Repo: "mikefarah/yq", Description: "YAML processor",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{".tar.gz", ".zip"},
					ArchiveFormat:     "none",
				},
			},
			{
				Name: "fzf", BinaryName: "fzf", Kind: GitHubRelease, Extension: true,
				Repo: "junegunn/fzf", Description: "Fuzzy finder",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tar.gz",
					BinaryPathInArchive: "fzf",
				},
			},
			{
				Name: "gh", BinaryName: "gh", Kind: GitHubRelease, Extension: true,
				Repo: "cli/cli", Description: "GitHub CLI",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					OSArchiveFormats:    map[string]string{"linux": "tar.gz", "darwin": "zip"},
					BinaryPathInArchive: "*/bin/gh",
				},
			},
			{
				Name: "tree-sitter", BinaryName: "tree-sitter", Kind: GitHubRelease, Extension: true,
				Repo: "tree-sitter/tree-sitter", Description: "Tree-sitter CLI (builds Neovim parsers)",
				Asset: AssetPattern{
					OSPatterns:         map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:       map[string]string{"amd64": "x64", "arm64": "arm64"},
					RequiredSubstrings: []string{"tree-sitter-cli-"},
					ArchiveFormat:      "zip",
				},
			},
			{
				Name: "gron", BinaryName: "gron", Kind: GitHubRelease, Extension: true,
				Repo: "tomnomnom/gron", Description: "Make JSON greppable",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "tgz",
					BinaryPathInArchive: "gron",
				},
			},
			{
				Name: "zoxide", BinaryName: "zoxide", Kind: GitHubRelease, Extension: true,
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
				Name: "sd", BinaryName: "sd", Kind: GitHubRelease, Extension: true,
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
				Name: "starship", BinaryName: "starship", Kind: GitHubRelease, Extension: true,
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
				Name: "anbu", BinaryName: "anbu", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/anbu", Description: "Anbu tool",
				Asset: AssetPattern{
					OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat:       "zip",
					BinaryPathInArchive: "anbu*",
				},
			},
			{
				Name: "danzo", BinaryName: "danzo", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/danzo", Description: "Danzo tool",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "starship-config", Kind: ConfigFile, Extension: true,
				Description: "Starship prompt configuration",
			},
			{
				Name: "age", BinaryName: "age", Kind: GitHubRelease, Extension: true,
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
				Name: "sq", BinaryName: "sq", Kind: GitHubRelease, Extension: true,
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
		Name:        "core",
		Description: "Dev tools, network utils, and media packages",
		Tools: []Tool{
			{
				Name: "dev-tools", Kind: SystemPackage, Extension: true,
				Description: "Development build tools",
				Platforms:   []string{"linux"},
				BrewPkgs:    []string{"cmake", "gcc", "make", "ninja", "gettext"},
			},
			{
				Name: "network-tools", Kind: SystemPackage, Extension: true,
				Description: "Network utilities",
				BrewPkgs:    []string{"nmap", "openssl"},
			},
			{
				Name: "media-tools", Kind: SystemPackage, Extension: true,
				Description: "Media and monitoring tools",
				BrewPkgs:    []string{"ffmpeg"},
			},
			{
				Name: "aerospace", Kind: SystemPackage, Extension: true,
				Description: "macOS tiling window manager",
				Platforms:   []string{"darwin"},
				BrewCasks:   []string{"nikitabobko/tap/aerospace"},
			},
		},
	},
	{
		Name:        "cloud",
		Description: "Cloud CLIs (AWS, Azure, GCP)",
		Tools: []Tool{
			{
				Name: "aws-cli", Kind: SystemPackage, Extension: true,
				Description: "AWS CLI v2",
				BrewPkgs:    []string{"awscli"},
			},
			{
				Name: "azure-cli", Kind: SystemPackage, Extension: true,
				Description: "Azure CLI",
				BrewPkgs:    []string{"azure-cli"},
			},
			{
				Name: "gcloud-cli", Kind: SystemPackage, Extension: true,
				Description: "Google Cloud CLI",
				BrewCasks:   []string{"gcloud-cli"},
			},
		},
	},
	{
		Name:        "runtimes",
		Description: "Language runtimes and language servers",
		Tools: []Tool{
			{
				Name: "uv", Kind: LanguageRuntime, Extension: true,
				Description: "Python package manager (binary only)",
			},
			{
				Name: "fnm", Kind: LanguageRuntime, Extension: true,
				Description: "Node version manager (binary only)",
			},
			{
				Name: "bun", BinaryName: "bun", Kind: GitHubRelease, Extension: true,
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
				Name: "go-sdk", Kind: LanguageRuntime, Extension: true,
				Description: "Go programming language SDK",
			},
			{
				Name: "java-sdk", Kind: LanguageRuntime, Extension: true,
				Description: "Eclipse Temurin JDK (latest LTS)",
			},
			{
				Name: "python", Kind: LanguageRuntime, Extension: true,
				Description: "Python via uv + py-default venv",
			},
			{
				Name: "rust", Kind: LanguageRuntime, Extension: true,
				Description: "Rust toolchain via rustup",
			},
			{
				Name: "node", Kind: LanguageRuntime, Extension: true,
				Description: "Node.js LTS via fnm",
			},
			{
				Name: "gopls", Kind: CustomScript, Extension: true,
				Description: "Go language server",
				InstallCmd:  "go install golang.org/x/tools/gopls@latest",
			},
			{
				Name: "pyright", Kind: NodePackage, Extension: true,
				Description: "Python language server (type checking)",
				NodePkg:     "pyright",
			},
			{
				Name: "typescript-language-server", Kind: NodePackage, Extension: true,
				Description: "TypeScript/JavaScript language server",
				// typescript@7 dropped the tsserver binary this server drives, so the major is pinned rather than tracking latest.
				NodePkg: "typescript-language-server typescript@5",
			},
			{
				Name: "ruff", Kind: PythonTool, Extension: true,
				Description: "Python linter and formatter with a language server",
				PyTool:      "ruff",
			},
		},
	},
	{
		Name:        "ai-tools",
		Description: "AI coding agents and LLM CLIs",
		Requires:    "runtimes",
		Tools: []Tool{
			{
				Name: "claude-code", Kind: NodePackage, Extension: true,
				Description: "Anthropic Claude Code CLI agent",
				NodePkg:     "@anthropic-ai/claude-code",
			},
			{
				Name: "antigravity", Kind: CustomScript, Extension: true,
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
				Name: "codex", Kind: NodePackage, Extension: true,
				Description: "OpenAI Codex CLI agent",
				NodePkg:     "@openai/codex",
			},
			{
				Name: "cursor-agent", Kind: CustomScript, Extension: true,
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
		Requires:    "runtimes",
		Tools: []Tool{
			{
				Name: "caddy", BinaryName: "caddy", Kind: GitHubRelease, Extension: true,
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
				Name: "linksnapper", BinaryName: "linksnapper", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/linksnapper", Description: "LinkSnapper bookmark manager",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "kairo", BinaryName: "kairo", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/kairo", Description: "Kairo markdown note-taking app",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "raikiri", BinaryName: "raikiri", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/raikiri", Description: "Self-hosted media and music server",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "expenseowl", BinaryName: "expenseowl", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/expenseowl", Description: "Self-hosted expense tracker",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "telly", BinaryName: "telly", Kind: GitHubRelease, Extension: true,
				Repo: "Tanq16/telly", Description: "Private Telly tool", IsPrivate: true,
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "yt-dlp", BinaryName: "yt-dlp", Kind: GitHubRelease, Extension: true,
				Repo: "yt-dlp/yt-dlp", Description: "Media downloader",
				Asset: AssetPattern{
					AssetNames: map[string]string{
						"linux/amd64":  "yt-dlp_linux",
						"linux/arm64":  "yt-dlp_linux_aarch64",
						"darwin/amd64": "yt-dlp_macos",
						"darwin/arm64": "yt-dlp_macos",
					},
					ArchiveFormat: "none",
				},
			},
			{
				Name: "rinnegan", Kind: AppBundle, Extension: true,
				Repo: "Tanq16/rinnegan", Description: "Self-contained PTY web terminal",
				Asset: AssetPattern{
					OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
					ArchPatterns:  map[string]string{"amd64": "x64", "arm64": "arm64"},
					ArchiveFormat: "tar.gz",
				},
			},
			{
				Name: "code-server", Kind: AppBundle, Extension: true,
				Repo: "coder/code-server", Description: "VS Code in the browser",
				PostInstall: "code-server",
				Asset: AssetPattern{
					OSPatterns:        map[string]string{"linux": "linux", "darwin": "macos"},
					ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
					ExcludeSubstrings: []string{".deb", ".rpm"},
					ArchiveFormat:     "tar.gz",
				},
			},
			{
				Name: "neo4j", Kind: AppBundle, Extension: true,
				Description:  "Graph database with cypher-shell and browser UI",
				URL:          "https://dist.neo4j.org/neo4j-community-{version}-unix.tar.gz",
				StableURL:    "https://repo1.maven.org/maven2/org/neo4j/neo4j/maven-metadata.xml",
				VersionRegex: `<release>([^<]+)</release>`,
				PostInstall:  "neo4j",
				Asset:        AssetPattern{ArchiveFormat: "tar.gz"},
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

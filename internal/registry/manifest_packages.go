package registry

var groups = []Group{
	{Name: GroupBrew, Description: "Brew-managed system packages and cloud CLIs"},
	{Name: GroupMacOSDesktop, Description: "macOS desktop apps and fonts"},
	{Name: GroupRuntime, Description: "Language runtimes and language servers"},
	{Name: GroupAI, Description: "AI coding agents and LLM CLIs"},
	{Name: GroupBinaries, Description: "Security, cloud and infrastructure binaries"},
	{Name: GroupHomelab, Description: "Self-hosted homelab services"},
	{Name: GroupPrivate, Description: "Personal tools"},
}

var suites = []Suite{
	{Name: "go-suite", Description: "Go SDK and gopls"},
	{Name: "js-suite", Description: "fnm, Node, bun and the JS/Python language servers"},
	{Name: "python-suite", Description: "uv, Python and ruff"},
	{Name: "java-suite", Description: "Eclipse Temurin JDK"},
	{Name: "rust-suite", Description: "Rust toolchain via rustup"},
}

var packages = []Tool{
	{
		Name: "media-tools", Kind: SystemPackage, Group: GroupBrew,
		Description: "Media processing tools",
		BrewPkgs:    []string{"ffmpeg", "imagemagick"},
	},
	{
		Name: "network-tools", Kind: SystemPackage, Group: GroupBrew,
		Description: "Network utilities",
		BrewPkgs:    []string{"nmap", "openssl"},
	},
	{
		Name: "aws-cli", Kind: SystemPackage, Group: GroupBrew,
		Description: "AWS CLI v2",
		BrewPkgs:    []string{"awscli"},
	},
	{
		Name: "azure-cli", Kind: SystemPackage, Group: GroupBrew,
		Description: "Azure CLI",
		BrewPkgs:    []string{"azure-cli"},
	},
	{
		Name: "gcloud-cli", Kind: SystemPackage, Group: GroupBrew,
		Description: "Google Cloud CLI",
		BrewCasks:   []string{"gcloud-cli"},
	},

	{
		Name: "aerospace", Kind: SystemPackage, Group: GroupMacOSDesktop,
		Description: "macOS tiling window manager",
		Platforms:   []string{"darwin"},
		BrewCasks:   []string{"nikitabobko/tap/aerospace"},
	},
	{
		Name: "sol", Kind: SystemPackage, Group: GroupMacOSDesktop,
		Description: "macOS launcher and command palette",
		Platforms:   []string{"darwin"},
		BrewCasks:   []string{"sol"},
	},
	{
		Name: "nerd-font", Kind: SystemPackage, Group: GroupMacOSDesktop,
		Description: "JetBrains Mono Nerd Font",
		Platforms:   []string{"darwin"},
		BrewCasks:   []string{"font-jetbrains-mono-nerd-font"},
	},

	{
		Name: "go-sdk", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "go-suite",
		Description: "Go programming language SDK",
	},
	{
		Name: "gopls", Kind: CustomScript, Group: GroupRuntime, Suite: "go-suite",
		Description: "Go language server",
		Requires:    []string{"go-sdk"},
		InstallCmd:  "go install golang.org/x/tools/gopls@latest",
	},
	{
		Name: "fnm", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "js-suite",
		Description: "Node version manager (binary only)",
	},
	{
		Name: "node", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "js-suite",
		Description: "Node.js LTS via fnm",
	},
	{
		Name: "bun", BinaryName: "bun", Kind: GitHubRelease, Group: GroupRuntime, Suite: "js-suite",
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
		Name: "pyright", Kind: NodePackage, Group: GroupRuntime, Suite: "js-suite",
		Description: "Python language server (type checking)",
		Requires:    []string{"node"},
		NodePkg:     "pyright",
	},
	{
		Name: "typescript-language-server", Kind: NodePackage, Group: GroupRuntime, Suite: "js-suite",
		Description: "TypeScript/JavaScript language server",
		Requires:    []string{"node"},
		// typescript@7 dropped the tsserver binary this server drives, so the major is pinned rather than tracking latest.
		NodePkg: "typescript-language-server typescript@5",
	},
	{
		Name: "uv", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "python-suite",
		Description: "Python package manager (binary only)",
	},
	{
		Name: "python", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "python-suite",
		Description: "Python via uv + py-default venv",
	},
	{
		Name: "ruff", Kind: PythonTool, Group: GroupRuntime, Suite: "python-suite",
		Description: "Python linter and formatter with a language server",
		Requires:    []string{"uv"},
		PyTool:      "ruff",
	},
	{
		Name: "java-sdk", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "java-suite",
		Description: "Eclipse Temurin JDK (latest LTS)",
	},
	{
		Name: "rust", Kind: LanguageRuntime, Group: GroupRuntime, Suite: "rust-suite",
		Description: "Rust toolchain via rustup",
	},

	{
		Name: "claude-code", Kind: NodePackage, Group: GroupAI,
		Description: "Anthropic Claude Code CLI agent",
		Requires:    []string{"node"},
		NodePkg:     "@anthropic-ai/claude-code",
	},
	{
		Name: "codex", Kind: NodePackage, Group: GroupAI,
		Description: "OpenAI Codex CLI agent",
		Requires:    []string{"node"},
		NodePkg:     "@openai/codex",
	},
	{
		Name: "antigravity", Kind: CustomScript, Group: GroupAI,
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
		Name: "cursor-agent", Kind: CustomScript, Group: GroupAI,
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

	{
		Name: "nuclei", BinaryName: "nuclei", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/nuclei", Description: "Vulnerability scanner",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "nuclei",
		},
	},
	{
		Name: "nuclei-templates", Kind: RepoSnapshot, Group: GroupBinaries,
		Description: "Nuclei vulnerability templates",
		Repo:        "projectdiscovery/nuclei-templates",
		Dest:        "~/shell/nuclei-templates",
	},
	{
		Name: "naabu", BinaryName: "naabu", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/naabu", Description: "Port scanner",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "naabu",
		},
	},
	{
		Name: "subfinder", BinaryName: "subfinder", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/subfinder", Description: "Subdomain discovery",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "subfinder",
		},
	},
	{
		Name: "proxify", BinaryName: "proxify", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/proxify", Description: "HTTP proxy",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "proxify",
		},
	},
	{
		Name: "httpx", BinaryName: "httpx", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/httpx", Description: "HTTP toolkit",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "httpx",
		},
	},
	{
		Name: "dnsx", BinaryName: "dnsx", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/dnsx", Description: "DNS toolkit",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "dnsx",
		},
	},
	{
		Name: "katana", BinaryName: "katana", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "projectdiscovery/katana", Description: "Web crawler",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "katana",
		},
	},
	{
		Name: "trufflehog", BinaryName: "trufflehog", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "trufflesecurity/trufflehog", Description: "Secret scanner",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "trufflehog",
		},
	},
	{
		Name: "ffuf", BinaryName: "ffuf", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "ffuf/ffuf", Description: "Fast web fuzzer",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "ffuf",
		},
	},
	{
		Name: "gobuster", BinaryName: "gobuster", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "OJ/gobuster", Description: "Directory/DNS brute-forcer",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "Linux", "darwin": "Darwin"},
			ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "gobuster",
		},
	},
	{
		Name: "gau", BinaryName: "gau", Kind: GitHubRelease, Group: GroupBinaries,
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
		Name: "gowitness", BinaryName: "gowitness", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "sensepost/gowitness", Description: "Web screenshot tool",
		Asset: AssetPattern{
			OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ExcludeSubstrings: []string{"windows"},
			ArchiveFormat:     "none",
		},
	},
	{
		Name: "kubelogin", BinaryName: "kubelogin", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "Azure/kubelogin", Description: "Azure Kubernetes login",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "bin/*/kubelogin",
		},
	},
	{
		Name: "grpcurl", BinaryName: "grpcurl", Kind: GitHubRelease, Group: GroupBinaries,
		Repo: "fullstorydev/grpcurl", Description: "curl for gRPC",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "osx"},
			ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "grpcurl",
		},
	},
	{
		Name: "terraform", BinaryName: "terraform", Kind: DirectDownload, Group: GroupBinaries,
		Repo: "hashicorp/terraform", Description: "Infrastructure as code",
		URL: "https://releases.hashicorp.com/terraform/{version_bare}/terraform_{version_bare}_{os}_{arch}.zip",
		Asset: AssetPattern{
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "terraform",
		},
	},
	{
		Name: "kubectl", BinaryName: "kubectl", Kind: DirectDownload, Group: GroupBinaries,
		Description: "Kubernetes CLI",
		StableURL:   "https://dl.k8s.io/release/stable.txt",
		URL:         "https://dl.k8s.io/release/{version}/bin/{os}/{arch}/kubectl",
	},
	{
		Name: "trivy", BinaryName: "trivy", Kind: GitHubRelease, Group: GroupBinaries,
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
		Name: "caddy", BinaryName: "caddy", Kind: GitHubRelease, Group: GroupHomelab,
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
		Name: "senkaimon", BinaryName: "senkaimon", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/senkaimon", Description: "Forward-auth identity service for a Caddy edge",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "linksnapper", BinaryName: "linksnapper", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/linksnapper", Description: "LinkSnapper bookmark manager",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "kairo", BinaryName: "kairo", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/kairo", Description: "Kairo markdown note-taking app",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "raikiri", BinaryName: "raikiri", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/raikiri", Description: "Self-hosted media and music server",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "expenseowl", BinaryName: "expenseowl", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/expenseowl", Description: "Self-hosted expense tracker",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "backhub", BinaryName: "backhub", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/backhub", Description: "GitHub repo mirror backups",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "zip",
		},
	},
	{
		Name: "local-content-share", BinaryName: "local-content-share", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/local-content-share", Description: "Self-hosted text and file sharing app",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "whiteboard", BinaryName: "whiteboard", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/whiteboard", Description: "Local-first infinite canvas whiteboard",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "inoichi", BinaryName: "inoichi", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/inoichi", Description: "Local-first mind mapping editor", IsPrivate: true,
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "goff", BinaryName: "goff", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/goff", Description: "TUI and CLI harness for FFmpeg",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "yt-dlp", BinaryName: "yt-dlp", Kind: GitHubRelease, Group: GroupHomelab,
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
		Name: "telly", BinaryName: "telly", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/telly", Description: "Telly tool", IsPrivate: true,
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "toon", BinaryName: "toon", Kind: GitHubRelease, Group: GroupHomelab,
		Repo: "Tanq16/toon", Description: "Private Toon tool", IsPrivate: true,
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "toon-*",
		},
	},
	{
		Name: "rinnegan", Kind: AppBundle, Group: GroupHomelab,
		Repo: "Tanq16/rinnegan", Description: "Self-contained PTY web terminal",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "x64", "arm64": "arm64"},
			ArchiveFormat: "tar.gz",
		},
	},
	{
		Name: "code-server", Kind: AppBundle, Group: GroupHomelab,
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
		Name: "neo4j", Kind: AppBundle, Group: GroupHomelab,
		Description:  "Graph database with cypher-shell and browser UI",
		URL:          "https://dist.neo4j.org/neo4j-community-{version}-unix.tar.gz",
		StableURL:    "https://repo1.maven.org/maven2/org/neo4j/neo4j/maven-metadata.xml",
		VersionRegex: `<release>([^<]+)</release>`,
		PostInstall:  "neo4j",
		Asset:        AssetPattern{ArchiveFormat: "tar.gz"},
	},

	{
		Name: "anbu", BinaryName: "anbu", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/anbu", Description: "Anbu tool",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "box", BinaryName: "box", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/box-cli", Description: "Box CLI tool",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "gcli", BinaryName: "gcli", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/gcli", Description: "Gcli tool",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "nits", BinaryName: "nits", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/nits", Description: "Nits tool",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "sharingan", BinaryName: "sharingan", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/sharingan", Description: "EC2 workstation manager",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "claudex", BinaryName: "claudex", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/claudex", Description: "Claudex tool",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "cybernest", BinaryName: "cybernest", Kind: GitHubRelease, Group: GroupPrivate,
		Repo: "Tanq16/cybernest", Description: "Private Cybernest tool", IsPrivate: true,
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
}

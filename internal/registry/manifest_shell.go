package registry

var shellTools = []Tool{
	{
		Name: "bat", BinaryName: "bat", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "fd", BinaryName: "fd", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "ripgrep", BinaryName: "rg", Kind: GitHubRelease, Group: GroupShell,
		Repo: "BurntSushi/ripgrep", Description: "Fast recursive grep",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
			ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "*/rg",
		},
	},
	{
		Name: "lsd", BinaryName: "lsd", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "jq", BinaryName: "jq", Kind: GitHubRelease, Group: GroupShell,
		Repo: "jqlang/jq", Description: "Command-line JSON processor",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "yq", BinaryName: "yq", Kind: GitHubRelease, Group: GroupShell,
		Repo: "mikefarah/yq", Description: "YAML processor",
		Asset: AssetPattern{
			OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ExcludeSubstrings: []string{".tar.gz", ".zip"},
			ArchiveFormat:     "none",
		},
	},
	{
		Name: "fzf", BinaryName: "fzf", Kind: GitHubRelease, Group: GroupShell,
		Repo: "junegunn/fzf", Description: "Fuzzy finder",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "fzf",
		},
	},
	{
		Name: "gh", BinaryName: "gh", Kind: GitHubRelease, Group: GroupShell,
		Repo: "cli/cli", Description: "GitHub CLI",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			OSArchiveFormats:    map[string]string{"linux": "tar.gz", "darwin": "zip"},
			BinaryPathInArchive: "*/bin/gh",
		},
	},
	{
		Name: "tree-sitter", BinaryName: "tree-sitter", Kind: GitHubRelease, Group: GroupShell,
		Repo: "tree-sitter/tree-sitter", Description: "Tree-sitter CLI (builds Neovim parsers)",
		Asset: AssetPattern{
			OSPatterns:         map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:       map[string]string{"amd64": "x64", "arm64": "arm64"},
			RequiredSubstrings: []string{"tree-sitter-cli-"},
			ArchiveFormat:      "zip",
		},
	},
	{
		Name: "gron", BinaryName: "gron", Kind: GitHubRelease, Group: GroupShell,
		Repo: "tomnomnom/gron", Description: "Make JSON greppable",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tgz",
			BinaryPathInArchive: "gron",
		},
	},
	{
		Name: "zoxide", BinaryName: "zoxide", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "sd", BinaryName: "sd", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "starship", BinaryName: "starship", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "age", BinaryName: "age", Kind: GitHubRelease, Group: GroupShell,
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
		Name: "sq", BinaryName: "sq", Kind: GitHubRelease, Group: GroupShell,
		Repo: "neilotoole/sq", Description: "jq-like data wrangler for databases",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ExcludeSubstrings:   []string{"windows", ".deb", ".rpm", ".apk", ".zst", "checksums"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "sq",
		},
	},

	{
		Name: "neovim", Kind: AppBundle, Group: GroupShell,
		Repo: "neovim/neovim", Description: "Neovim text editor",
		PostInstall: "neovim",
		Asset: AssetPattern{
			OSPatterns:        map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:      map[string]string{"amd64": "x86_64", "arm64": "arm64"},
			ExcludeSubstrings: []string{".appimage", ".zsync", ".msi", ".zip"},
			ArchiveFormat:     "tar.gz",
		},
	},

	{
		Name: "zsh-autosuggestions", Kind: RepoSnapshot, Group: GroupShell,
		Description: "ZSH autosuggestions plugin",
		Repo:        "zsh-users/zsh-autosuggestions",
		Dest:        "~/shell/plugins/zsh-autosuggestions",
	},
	{
		Name: "zsh-syntax-highlighting", Kind: RepoSnapshot, Group: GroupShell,
		Description: "ZSH syntax highlighting plugin",
		Repo:        "zsh-users/zsh-syntax-highlighting",
		Dest:        "~/shell/plugins/zsh-syntax-highlighting",
	},

	{
		Name: "tmux-config", Kind: ConfigFile, Group: GroupShell,
		Description: "Tmux configuration",
	},
	{
		Name: "kitty-config", Kind: ConfigFile, Group: GroupShell,
		Description: "Kitty terminal configuration",
	},
	{
		Name: "kitty-theme", Kind: ConfigFile, Group: GroupShell,
		Description: "Kitty theme configuration",
	},
	{
		Name: "lsd-colors", Kind: ConfigFile, Group: GroupShell,
		Description: "lsd colour theme",
	},
	{
		Name: "aerospace-config", Kind: ConfigFile, Group: GroupShell,
		Description: "Aerospace WM configuration",
		Platforms:   []string{"darwin"},
	},
	{
		Name: "rcfile", Kind: ConfigFile, Group: GroupShell,
		Description: "Zsh RC file (complete .zshrc)",
	},
	{
		Name: "nvim-config", Kind: ConfigFile, Group: GroupShell,
		Description: "Neovim configuration",
	},
	{
		Name: "starship-config", Kind: ConfigFile, Group: GroupShell,
		Description: "Starship prompt configuration",
	},
}

package registry

var AllTools = []Tool{
	// ========== GitHub Release Binaries (Public) ==========
	{
		Name: "bat", BinaryName: "bat", Kind: GitHubRelease, Category: Core,
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
		Name: "fd", BinaryName: "fd", Kind: GitHubRelease, Category: Core,
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
		Name: "ripgrep", BinaryName: "rg", Kind: GitHubRelease, Category: Core,
		Repo: "BurntSushi/ripgrep", Description: "Fast recursive grep",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
			ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "*/rg",
		},
	},
	{
		Name: "lsd", BinaryName: "lsd", Kind: GitHubRelease, Category: Core,
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
		Name: "jq", BinaryName: "jq", Kind: GitHubRelease, Category: Core,
		Repo: "jqlang/jq", Description: "Command-line JSON processor",
		Asset: AssetPattern{
			OSPatterns:    map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:  map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat: "none",
		},
	},
	{
		Name: "yq", BinaryName: "yq", Kind: GitHubRelease, Category: Core,
		Repo: "mikefarah/yq", Description: "YAML processor",
		Asset: AssetPattern{
			OSPatterns:        map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:      map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ExcludeSubstrings: []string{".tar.gz", ".zip"},
			ArchiveFormat:     "none",
		},
	},
	{
		Name: "fzf", BinaryName: "fzf", Kind: GitHubRelease, Category: Core,
		Repo: "junegunn/fzf", Description: "Fuzzy finder",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "fzf",
		},
	},
	{
		Name: "gh", BinaryName: "gh", Kind: GitHubRelease, Category: Core,
		Repo: "cli/cli", Description: "GitHub CLI",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macOS"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			OSArchiveFormats:    map[string]string{"linux": "tar.gz", "darwin": "zip"},
			BinaryPathInArchive: "*/bin/gh",
		},
	},
	{
		Name: "gron", BinaryName: "gron", Kind: GitHubRelease, Category: Core,
		Repo: "tomnomnom/gron", Description: "Make JSON greppable",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tgz",
			BinaryPathInArchive: "gron",
		},
	},
	{
		Name: "sq", BinaryName: "sq", Kind: GitHubRelease, Category: Core,
		Repo: "neilotoole/sq", Description: "Data wrangler",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "macos"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "sq",
		},
	},
	{
		Name: "zoxide", BinaryName: "zoxide", Kind: GitHubRelease, Category: Core,
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
		Name: "sd", BinaryName: "sd", Kind: GitHubRelease, Category: Core,
		Repo: "chmln/sd", Description: "Find and replace CLI tool",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "apple"},
			ArchPatterns:        map[string]string{"amd64": "x86_64", "arm64": "aarch64"},
			ExcludeSubstrings:   []string{"windows", "gnueabi"},
			ArchiveFormat:       "tar.gz",
			BinaryPathInArchive: "*/sd",
		},
	},

	// ========== Own Public Tools (Tanq16) ==========
	{
		Name: "anbu", BinaryName: "anbu", Kind: GitHubRelease, Category: Core,
		Repo: "Tanq16/anbu", Description: "Anbu tool",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "anbu*",
		},
	},
	{
		Name: "danzo", BinaryName: "danzo", Kind: GitHubRelease, Category: Core,
		Repo: "Tanq16/danzo", Description: "Danzo tool",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "danzo",
		},
	},
	{
		Name: "ai-context", BinaryName: "ai-context", Kind: GitHubRelease, Category: Core,
		Repo: "Tanq16/ai-context", Description: "AI context builder",
		Asset: AssetPattern{
			OSPatterns:          map[string]string{"linux": "linux", "darwin": "darwin"},
			ArchPatterns:        map[string]string{"amd64": "amd64", "arm64": "arm64"},
			ArchiveFormat:       "zip",
			BinaryPathInArchive: "ai-context",
		},
	},

	// ========== System Packages (base) ==========
	{
		Name: "core-utils", Kind: SystemPackage, Category: System,
		Description: "Core system utilities",
		BrewPkgs:    []string{"wget", "zip", "unzip", "file"},
	},
	{
		Name: "shell-base", Kind: SystemPackage, Category: System,
		Description: "Shell and terminal essentials",
		BrewPkgs:    []string{"tmux", "htop"},
	},

	// ========== Neovim (via brew) ==========
	{
		Name: "neovim", Kind: SystemPackage, Category: System,
		Description: "Neovim text editor (0.11+ for NvChad)",
		BrewPkgs:    []string{"neovim"},
	},

	// ========== Config Files ==========
	{
		Name: "tmux-config", Kind: ConfigFile, Category: Config,
		Description: "Tmux configuration",
	},
	{
		Name: "kitty-config", Kind: ConfigFile, Category: Config,
		Description: "Kitty terminal configuration",
	},
	{
		Name: "kitty-theme", Kind: ConfigFile, Category: Config,
		Description: "Kitty theme configuration",
	},
	{
		Name: "aerospace-config", Kind: ConfigFile, Category: Config,
		Description: "Aerospace WM configuration",
		Platforms:   []string{"darwin"},
	},
	{
		Name: "rcfile", Kind: ConfigFile, Category: Config,
		Description: "Zsh RC file (complete .zshrc)",
	},

	// ========== Shell Plugins ==========
	{
		Name: "spaceship-prompt", Kind: ShellPlugin, Category: Shell,
		Description: "Spaceship ZSH theme",
		CloneURL:    "https://github.com/spaceship-prompt/spaceship-prompt.git",
		CloneDest:   "~/.oh-my-zsh/custom/themes/spaceship-prompt",
		PostClone:   "spaceship",
	},
	{
		Name: "zsh-autosuggestions", Kind: ShellPlugin, Category: Shell,
		Description: "ZSH autosuggestions plugin",
		CloneURL:    "https://github.com/zsh-users/zsh-autosuggestions.git",
		CloneDest:   "~/.oh-my-zsh/custom/plugins/zsh-autosuggestions",
	},
	{
		Name: "zsh-syntax-highlighting", Kind: ShellPlugin, Category: Shell,
		Description: "ZSH syntax highlighting plugin",
		CloneURL:    "https://github.com/zsh-users/zsh-syntax-highlighting.git",
		CloneDest:   "~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting",
	},
	{
		Name: "tpm", Kind: ShellPlugin, Category: Shell,
		Description: "Tmux Plugin Manager",
		CloneURL:    "https://github.com/tmux-plugins/tpm.git",
		CloneDest:   "~/.tmux/plugins/tpm",
		PostClone:   "tpm",
	},
	{
		Name: "nvchad", Kind: ShellPlugin, Category: Shell,
		Description: "NvChad Neovim configuration",
		CloneURL:    "https://github.com/NvChad/starter.git",
		CloneDest:   "~/.config/nvim",
		PostClone:   "nvchad",
	},
}

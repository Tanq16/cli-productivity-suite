package registry

var AllTools = []Tool{
	{
		Name: "core-utils", Kind: SystemPackage,
		Description: "Core system utilities",
		BrewPkgs:    []string{"wget", "zip", "unzip", "file"},
	},
	{
		Name: "shell-base", Kind: SystemPackage,
		Description: "Shell and terminal essentials",
		BrewPkgs:    []string{"tmux", "htop"},
	},

	{
		Name: "neovim", Kind: AppBundle,
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
		Name: "tmux-config", Kind: ConfigFile,
		Description: "Tmux configuration",
	},
	{
		Name: "kitty-config", Kind: ConfigFile,
		Description: "Kitty terminal configuration",
	},
	{
		Name: "kitty-theme", Kind: ConfigFile,
		Description: "Kitty theme configuration",
	},
	{
		Name: "lsd-colors", Kind: ConfigFile,
		Description: "lsd colour theme",
	},
	{
		Name: "aerospace-config", Kind: ConfigFile,
		Description: "Aerospace WM configuration",
		Platforms:   []string{"darwin"},
	},
	{
		Name: "rcfile", Kind: ConfigFile,
		Description: "Zsh RC file (complete .zshrc)",
	},
	{
		Name: "nvim-config", Kind: ConfigFile,
		Description: "Neovim configuration",
	},

	{
		Name: "zsh-autosuggestions", Kind: RepoSnapshot,
		Description: "ZSH autosuggestions plugin",
		Repo:        "zsh-users/zsh-autosuggestions",
		Dest:        "~/shell/plugins/zsh-autosuggestions",
	},
	{
		Name: "zsh-syntax-highlighting", Kind: RepoSnapshot,
		Description: "ZSH syntax highlighting plugin",
		Repo:        "zsh-users/zsh-syntax-highlighting",
		Dest:        "~/shell/plugins/zsh-syntax-highlighting",
	},
}

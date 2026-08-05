package registry

var AllTools = []Tool{
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

	{
		Name: "neovim", Kind: AppBundle, Category: System,
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
	{
		Name: "nvim-config", Kind: ConfigFile, Category: Config,
		Description: "Neovim configuration",
	},

	{
		Name: "zsh-autosuggestions", Kind: ShellPlugin, Category: Shell,
		Description: "ZSH autosuggestions plugin",
		CloneURL:    "https://github.com/zsh-users/zsh-autosuggestions.git",
		CloneDest:   "~/shell/plugins/zsh-autosuggestions",
	},
	{
		Name: "zsh-syntax-highlighting", Kind: ShellPlugin, Category: Shell,
		Description: "ZSH syntax highlighting plugin",
		CloneURL:    "https://github.com/zsh-users/zsh-syntax-highlighting.git",
		CloneDest:   "~/shell/plugins/zsh-syntax-highlighting",
	},
}

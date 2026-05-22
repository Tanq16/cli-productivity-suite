package registry

var AllTools = []Tool{
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

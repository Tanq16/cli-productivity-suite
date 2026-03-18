<div align="center">
  <img src=".github/assets/logo.png" alt="CLI Productivity Suite Logo" width="200">
  <h1>CLI Productivity Suite</h1>

  <a href="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/cli-productivity-suite/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/cli-productivity-suite"></a><br><br>
  <a href="#capabilities">Capabilities</a> &bull; <a href="#installation">Installation</a> &bull; <a href="#usage">Usage</a> &bull; <a href="#tips-and-notes">Tips & Notes</a> &bull; <a href="#deep-removal">Deep Removal</a>
</div>

---

A single Go binary (`cps`) to initialize, manage, and update a complete CLI-driven development environment on Linux and macOS. It installs and tracks ~50 tools across these categories:

- **CLI utilities** — bat, fd, ripgrep, lsd, jq, yq, fzf, gron, sq, and more
- **Security tools** — nuclei, naabu, katana, subfinder, ffuf, gobuster, trufflehog, proxify
- **Cloud & infra** — AWS CLI, Azure CLI, gcloud CLI, terraform, kubectl, kubelogin
- **Language runtimes** — Go SDK, Python 3.14 (via uv), Node.js LTS (via nvm)
- **Editor & shell** — Neovim (0.11+) with NvChad, spaceship-prompt, zsh plugins, tmux with TPM
- **Config files** — complete `.zshrc`, tmux.conf, kitty.conf, aerospace.toml (macOS)

## Capabilities

| Category | Commands | Description |
|----------|----------|-------------|
| Setup | `cps init` | Full environment setup — system packages, ~50 tools, cloud CLIs, language runtimes, shell plugins, and config files |
| Monitoring | `cps check` | Compare installed versions against latest releases |
| Updates | `cps update`, `cps install <tool>` | Update all installed tools or install a single tool by name |
| Maintenance | `cps clean` | Remove all CPS-managed files and directories |

## Installation

### Binary

Download from [releases](https://github.com/tanq16/cli-productivity-suite/releases):

```bash
# Linux/macOS
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64; [ "$ARCH" = "aarch64" ] && ARCH=arm64
curl -sL https://github.com/tanq16/cli-productivity-suite/releases/latest/download/cps-$(uname -s | tr '[:upper:]' '[:lower:]')-$ARCH -o cps
chmod +x cps
sudo mv cps /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/tanq16/cli-productivity-suite
cd cli-productivity-suite
make build
```

### Prerequisites

- [Oh My Zsh](https://ohmyz.sh/) must be installed before running `cps init`
- Git must be available in PATH
- Install a nerd font for your terminal emulator — recommended: [JetBrains Mono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/JetBrainsMono.zip)

## Usage

### `init`

Full environment setup with all tools, plugins, and configs.

```bash
# Public tools only
cps init

# Include private repos
cps init --gh-token YOUR_GITHUB_PAT
```

### `check`

Compare installed tool versions against latest releases.

```bash
cps check              # All installed tools
cps check --public     # Public tools only
cps check --private    # Private tools only
cps check --system     # System packages only
```

### `update`

Update installed tools to their latest versions.

```bash
cps update             # Update all
cps update --public    # Public tools only
```

### `install`

Install or update a single tool by name.

```bash
cps install bat
cps install terraform
cps install --gh-token TOKEN gcli
```

### `clean`

Remove all CPS-managed files and directories with confirmation.

```bash
cps clean
```

**Global flags:**

| Flag | Description |
|------|-------------|
| `--debug` | Enable debug logging |
| `--for-ai` | AI-friendly output (markdown tables, no color) |
| `--gh-token` | GitHub PAT for private repos (falls back to `gh auth token`) |

## Tips and Notes

- All binary tools are installed to `~/shell/executables/` — add this to your PATH (the deployed `.zshrc` handles this automatically)
- Neovim is installed from GitHub releases (0.11+) on both Linux and macOS to meet NvChad requirements
- State is tracked in `~/.config/cps/state.json` — this file records installed versions for `check` and `update` commands
- If the `gh` CLI is authenticated (`gh auth login`), CPS automatically uses its token — no need to pass `--gh-token`
- Running `cps init` is idempotent — it skips tools that are already at the latest version
- Cloud CLIs (AWS, Azure, gcloud) require sudo on Linux for system-level installation
- The `.zshrc` deployed by `cps init` is a complete replacement — it includes Oh My Zsh config, all tool integrations, aliases, and functions
- `cps clean` removes `~/shell`, `~/.tmux`, `~/.config/nvim`, `~/.nvm`, `~/nuclei-templates`, `~/google-cloud-sdk`, and `~/.config/cps` — it does not touch Oh My Zsh, deployed config files, or system packages

## Deep Removal

The `cps clean` command performs a superficial cleanup of CPS-managed directories. For a full removal of everything CPS installs, follow these steps:

**Step 1** — Run `cps clean` to remove the primary managed directories (`~/shell`, `~/.tmux`, `~/.config/nvim`, `~/.nvm`, `~/nuclei-templates`, `~/google-cloud-sdk`, `~/.config/cps`):

```bash
cps clean
```

**Step 2** — Remove remaining configs, data directories, Oh My Zsh, and system-level installs:

```bash
rm -rf \
  $HOME/.oh-my-zsh \
  $HOME/.tmux.conf \
  $HOME/.zshrc \
  $HOME/.aerospace.toml \
  $HOME/.config/kitty/kitty.conf \
  $HOME/.local/share/nvim \
  && sudo rm -rf /usr/local/go /usr/local/aws-cli /usr/local/bin/cps
```

**Step 3** — Remove system packages installed by CPS:

Linux (apt):

```bash
sudo apt-get remove -y tmux git tree wget curl zsh openssl nmap ncat cmake gcc make ninja-build gettext zip unzip file ffmpeg
```

macOS (brew):

```bash
brew uninstall tmux git tree wget curl openssl nmap ffmpeg
brew uninstall --cask nikitabobko/tap/aerospace
```

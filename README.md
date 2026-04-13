<div align="center">
  <img src=".github/assets/logo.png" alt="CLI Productivity Suite Logo" width="200">
  <h1>CLI Productivity Suite</h1>

  <a href="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/cli-productivity-suite/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/cli-productivity-suite"></a>
</div>

---

A single Go binary (`cps`) that sets up and manages a complete CLI development environment on **Linux** and **macOS**. Run `cps init` once to get a working shell with core tools, Neovim, tmux, and configs. Extend it with `cps extend` for language runtimes, cloud CLIs, security tools, and more.

## Prerequisites

| Requirement | Why |
|---|---|
| [Oh My Zsh](https://ohmyz.sh/) | Shell framework — `cps init` will not run without it |
| Git | Used to clone plugins and configs |
| [Homebrew](https://brew.sh/) | System package installs (both Linux and macOS) — `cps` installs all system/cloud CLI packages via brew |

**Recommended:**

- [Kitty](https://sw.kovidgoez.net/kitty/) terminal — CPS deploys a Kitty config and Catppuccin theme. Without Kitty, those config files are harmless but unused.
- [JetBrains Mono Nerd Font](https://www.nerdfonts.com/font-downloads) — the Kitty and Neovim configs expect a nerd font. Without one, icons and glyphs will render as boxes.

## Install

```bash
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64; [ "$ARCH" = "aarch64" ] && ARCH=arm64
mkdir -p "$HOME/.local/bin"
curl -sL "https://github.com/tanq16/cli-productivity-suite/releases/latest/download/cps-$(uname -s | tr '[:upper:]' '[:lower:]')-$ARCH" -o "$HOME/.local/bin/cps"
chmod +x "$HOME/.local/bin/cps"
```

If `~/.local/bin` isn't on your PATH yet (common on fresh macOS), run `cps init` via its full path — `~/.local/bin/cps init` — for the first invocation. The rc fragment that `init` deploys adds `~/.local/bin` to PATH for all future sessions.

Or build from source:

```bash
git clone https://github.com/tanq16/cli-productivity-suite && cd cli-productivity-suite
make build   # produces ./cps
```

## Usage

### `cps init`

Sets up the base environment — core CLI tools (bat, fd, ripgrep, lsd, jq, yq, fzf, sd, gron, sq, zoxide, gh, anbu, danzo, ai-context), Neovim with NvChad, zsh plugins, tmux with TPM, and config files. System packages install via Homebrew on both Linux and macOS — no sudo required.

```bash
cps init
cps extend core
```

These two commands are the standard way to set up this suite. `init` handles the base shell environment and `core` adds dev tools, network utilities, and media packages. Everything else via `cps extend` is optional — install what you need.

### `cps extend <pack> [tools...]`

Install extension packs or pick specific tools from a pack.

```bash
cps extend list                       # list all packs
cps extend runtimes                   # install all language runtimes
cps extend runtimes go-sdk            # install only Go
cps extend security nuclei subfinder  # pick specific tools
```

| Pack | Contents |
|---|---|
| core | Dev tools, network utils, media packages (cmake, nmap, ffmpeg, aerospace) |
| runtimes | uv, fnm, bun, Go, Python (includes uv), Rust, Node.js LTS (includes fnm) |
| cloud | AWS CLI, Azure CLI, gcloud CLI |
| security | nuclei, naabu, subfinder, proxify, httpx, dnsx, trufflehog, gobuster, titus, nuclei-templates |
| cloudsec | terraform, kubectl, kubelogin, grpcurl, cloudfox, aurelian, trivy, cloudlist |
| appsec | katana, ffuf, hadrian, dalfox, reaper, poltergeist, wraith, gau |
| misc | julius, trajan, gowitness, snitch, age |
| private | Personal tools (requires `--gh-token`) |

Packs with shell integration (`runtimes`, `cloud`, `security`) deploy RC fragments automatically.

### `cps cheat <topic>`

Terminal cheat sheets — `cps`, `go`, `uv`, `fnm`, `rust`, `tmux`, `nvim`, `fzf`, `regex`.

### `cps self-update`

Updates the `cps` binary in place (at whatever path it's running from).

### Flags

| Flag | Description |
|---|---|
| `--gh-token` | GitHub PAT for private repos (falls back to `gh auth token`) |
| `--debug` | Verbose debug logging |
| `--for-ai` | AI-friendly output (no color) |

## Custom Extension Packs

Drop a YAML file in `~/.config/cps/extensions/` to define your own pack:

```yaml
name: my-tools
description: My custom tools
shell:
  env:
    MY_VAR: "value"
  path_prepend:
    - "$HOME/.local/bin"
  source:
    - "$HOME/.cargo/env"
tools:
  - name: my-tool
    install: curl -sL https://example.com/install.sh | bash
```

Then run `cps extend my-tools`. Custom packs appear in `cps extend list` alongside built-in packs.

The `shell` block controls what gets added to your shell environment via a generated RC fragment at `~/shell/rc/custom/<pack-name>.zsh`:

- **`env`** — key-value pairs exported as environment variables
- **`path_prepend`** — directories prepended to `$PATH`
- **`source`** — files conditionally sourced (only if they exist)

All three fields are optional. If the entire `shell` block is omitted, no fragment is generated.

## Shell Integration

CPS uses a modular fragment system instead of a monolithic `.zshrc`:

| Fragment | Deployed by |
|---|---|
| `~/shell/rc/00-base.zsh` | `cps init` |
| `~/shell/rc/10-runtimes.zsh` | `cps extend runtimes` |
| `~/shell/rc/20-cloud.zsh` | `cps extend cloud` |
| `~/shell/rc/30-security.zsh` | `cps extend security` |
| `~/shell/rc/custom/*.zsh` | Custom packs or user-managed |

`~/.zshrc` is a thin loader that sources all fragments in order.

## Notes

- Core tools install to `~/shell/executables/`, extensions to `~/shell/extensions/` — both on PATH
- State tracked in `~/.config/cps/state.json` — runs are idempotent, already-current tools are skipped
- If `gh` CLI is authenticated, CPS uses its token automatically — no need for `--gh-token`

## Deep Removal

```bash
# CPS directories
rm -rf $HOME/shell $HOME/.tmux $HOME/.config/nvim $HOME/.config/cps

# Configs
rm -f $HOME/.tmux.conf $HOME/.zshrc $HOME/.aerospace.toml $HOME/.config/kitty/kitty.conf $HOME/.config/kitty/current-theme.conf
rm -rf $HOME/.local/share/nvim $HOME/.oh-my-zsh
rm -f $HOME/.local/bin/cps

# Uninstall brew-managed system and cloud packages (optional)
brew uninstall awscli azure-cli gcloud-cli
```

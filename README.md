<div align="center">
  <img src=".github/assets/logo.png" alt="CLI Productivity Suite Logo" width="200">
  <h1>CLI Productivity Suite</h1>

  <a href="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/cli-productivity-suite/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/cli-productivity-suite"></a>
</div>

---

A single Go binary (`cps`) that sets up and manages a complete CLI development environment on **Linux** and **macOS**. Run `cps init` once to get a working shell with core tools, Neovim, tmux, and configs. Extend it with `cps extend` for language runtimes, cloud CLIs, security tools, and more.

## Prerequisites

### Bootstrap

macOS ships zsh and curl out of the box; `git` comes from the Xcode Command Line Tools, so install those first if you haven't already:

```bash
xcode-select --install
```

On Linux:

```bash
sudo apt install git curl zsh build-essential
```

### Install these (both platforms)

| Requirement | One-line install |
|---|---|
| [Oh My Zsh](https://ohmyz.sh/) | `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"` |
| [Homebrew](https://brew.sh/) | `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"` |

Both are required — `cps init` won't run without them. CPS uses brew for all system and cloud CLI packages.

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

Sets up the base shell environment — Homebrew packages (`wget`, `zip`, `unzip`, `file`, `tmux`, `htop`, `neovim`), Neovim with NvChad, zsh plugins (autosuggestions, syntax-highlighting), tmux with TPM, and CPS-managed config files (`.zshrc`, `.tmux.conf`, kitty configs). No CLI binaries are installed here — those live in the `essentials` pack so they can be updated individually. No sudo required.

```bash
cps init
cps extend essentials
cps extend core
```

These three commands are the standard way to set up the suite. `init` handles the base shell environment, `essentials` adds the everyday CLI binaries (bat, fd, ripgrep, fzf, starship, etc.) and deploys the starship prompt config, and `core` adds dev tools, network utilities, and media packages. Everything else via `cps extend` is optional — install what you need.

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
| essentials | Everyday CLI binaries (bat, fd, ripgrep, lsd, jq, yq, fzf, gh, gron, zoxide, sd, starship, anbu, danzo, ai-context) + starship.toml |
| core | Dev tools, network utils, media packages (cmake, nmap, ffmpeg, aerospace) |
| runtimes | uv, fnm, bun, Go, Java (Temurin LTS), Python (includes uv), Rust, Node.js LTS (includes fnm) |
| cloud | AWS CLI, Azure CLI, gcloud CLI |
| security | nuclei, naabu, subfinder, proxify, httpx, dnsx, trufflehog, gobuster, nuclei-templates |
| cloudsec | terraform, kubectl, kubelogin, grpcurl, cloudfox, trivy, cloudlist |
| appsec | katana, ffuf, dalfox, reaper, poltergeist, wraith, gau |
| misc | gowitness, snitch, age |
| private | Personal tools (requires `--gh-token`) |

Packs with shell integration (`runtimes`, `cloud`, `security`) deploy RC fragments automatically.

### `cps cheat <topic>`

Terminal cheat sheets — `cps`, `go`, `java`, `uv`, `fnm`, `rust`, `tmux`, `nvim`, `fzf`, `regex`.

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

The repo's `custom-extensions/` directory ships ready-made reference packs (`ai-tools`, `additional-cloud-tools`, `database`, `praetorian`) — copy any of them to `~/.config/cps/extensions/` to use as-is, or treat them as templates for your own. Or pull all of them in one shot with `cps download-known-extensions` (below).

### `cps download-known-extensions`

Fetches the reference custom-extension YAMLs maintained in the CPS repo (`ai-tools`, `additional-cloud-tools`, `database`, `praetorian`) and writes them to `~/.config/cps/extensions/`. After running, they show up in `cps extend list` and you can install any of them with `cps extend <pack>` (or `cps extend <pack> <tool>` for a single tool).

```bash
cps download-known-extensions
cps extend list                  # ai-tools, database, etc. now visible
cps extend ai-tools claude-code  # install just claude-code from ai-tools
```

Overwrites existing files of the same name — if you've customized one of the reference packs locally, rename it before re-running.

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

Run the included script to wipe CPS, CPS-installed brew packages, and Oh My Zsh. Homebrew itself and `~/.zsh_history` are preserved so you can reinstall cleanly without rebuilding your shell history or re-bootstrapping brew.

```bash
./deep-removal.sh
```

<div align="center">
  <img src=".github/assets/logo.png" alt="CLI Productivity Suite Logo" width="200">
  <h1>CLI Productivity Suite</h1>

  <a href="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/cli-productivity-suite/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/cli-productivity-suite"></a><br><br>

  <a href="#prerequisites">Prerequisites</a> &bull; <a href="#install">Install</a> &bull; <a href="#usage">Usage</a> &bull; <a href="#custom-extension-packs">Custom Extensions</a> &bull; <a href="#shell-integration">Shell Integration</a> &bull; <a href="#sandbox-container">Sandbox Container</a> &bull; <a href="#deep-removal">Deep Removal</a>
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

### Install this (both platforms)

| Requirement | One-line install |
|---|---|
| [Homebrew](https://brew.sh/) | `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"` |

Homebrew is required — `cps init` won't run without it. CPS uses brew for all system and cloud CLI packages. CPS does **not** use Oh My Zsh; the two zsh plugins it installs (`zsh-autosuggestions`, `zsh-syntax-highlighting`) are sourced directly from `~/shell/plugins/`.

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

### Standard setup — run these three commands in order

```bash
cps init                  # base shell env, zsh plugins, configs, neovim, tmux
cps extend essentials     # everyday CLI binaries (bat, fd, ripgrep, fzf, starship, ...)
cps extend core           # dev/network/media brew packages (cmake, nmap, ffmpeg, aerospace)
```

> **Order matters.** `init` deploys `~/.zshrc` and rc fragments that reference binaries (`lsd`, `fd`, `bat`, `fzf`, `starship`, etc.) which the `essentials` pack provides. The shell still works without `essentials`, but aliases like `tree`/`l`/`la` and the starship prompt won't resolve until it's installed. Always follow `init` with `cps extend essentials`.

### `cps init`

Sets up the base shell environment — Homebrew packages (`wget`, `zip`, `unzip`, `file`, `tmux`, `htop`, `neovim`), Neovim with NvChad, zsh plugins (autosuggestions, syntax-highlighting), tmux with TPM, and CPS-managed config files (`.zshrc`, `.tmux.conf`, kitty configs). No CLI binaries are installed here — those live in the `essentials` pack so they can be updated individually. No sudo required.

Everything else via `cps extend` is optional — install what you need.

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

Terminal cheat sheets — `cps`, `go`, `java`, `uv`, `fnm`, `bun`, `rust`, `tmux`, `nvim`, `fzf`, `jq`, `regex`.

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

### Adding your own stuff (no extension pack needed)

Three drop-zones, three buckets — you never need to write a YAML pack for personal tweaks:

| You want to add… | Drop it here | Notes |
|---|---|---|
| Aliases, exports, functions, custom sourcing | `~/shell/rc/custom/anything.zsh` | Loaded automatically by `~/.zshrc` after CPS fragments — so it can override CPS-set values |
| Your own binaries / scripts | `~/shell/custom-bin/` | Prepended to PATH **ahead of** CPS-managed dirs, so your binary wins if a name collides with a CPS one |
| A reusable, idempotent install bundle you want `cps extend` to manage | `~/.config/cps/extensions/<name>.yaml` | See [Custom Extension Packs](#custom-extension-packs) above |

Both `~/shell/rc/custom/` and `~/shell/custom-bin/` are created by `cps init` and are **never touched** by subsequent `cps init` / `cps extend` runs.

`deep-removal.sh` wipes the whole `~/shell/` tree, so anything you drop there is removed by it — if you want long-term-survival storage, keep it elsewhere.

## Notes

- Core tools install to `~/shell/executables/`, extensions to `~/shell/extensions/` — both on PATH
- User-owned binaries live in `~/shell/custom-bin/` (also on PATH, prepended ahead of the CPS-managed dirs)
- State tracked in `~/.config/cps/state.json` — runs are idempotent, already-current tools are skipped
- If `gh` CLI is authenticated, CPS uses its token automatically — no need for `--gh-token`
- `00-base.zsh` exports `HOMEBREW_NO_AUTO_UPDATE=1` so `brew install` stays fast and deterministic. If you want brew to auto-update on every invocation, drop `unset HOMEBREW_NO_AUTO_UPDATE` into a file under `~/shell/rc/custom/`

## Sandbox Container

A prebuilt Ubuntu container with the full CPS environment baked in — every built-in extension pack (except `private`), brew on Linuxbrew, a non-root `cps` user with sudo, and zsh + tmux + neovim ready to go. Useful when you need a CPS-style workspace on a machine where you can't (or don't want to) install CPS directly.

```bash
docker run -d --name cps-sandbox tanq16/cps-sandbox:latest
docker exec -it cps-sandbox zsh -l
```

The image runs `sleep infinity` as its default command, so it stays alive and you `docker exec` in whenever you need it. `docker exec -it <name> zsh -l` always gives you the full configured shell (rc fragments sourced, PATH wired up, starship prompt, plugins loaded). Inside the shell, `tt` starts a tmux session, `t` re-attaches.

Build locally:

```bash
docker build -t cps-sandbox .
docker run -d --name cps-sandbox cps-sandbox
docker exec -it cps-sandbox zsh -l
```

The image is multi-arch (`linux/amd64` + `linux/arm64`) and large (multi-GB) — it carries full language runtimes, cloud CLIs, security tooling, every reference custom-extension pack (`ai-tools`, `additional-cloud-tools`, `database`, `praetorian`), and the public-repo tools from the `private` pack (`nits`, `raikiri`, `gcli`, `box`, `claudex`). The four truly-private tools (`toon`, `nblm`, `cybernest`, `lincli`) are skipped since they need an auth token.

### A ready environment for AI agents

The prebuilt image is intentionally a **drop-in toolkit for AI coding agents** — Claude Code, Codex, opencode, Crush, antigravity, and friends. Spin up the container once and a single non-root user already has:

- **The agent CLIs themselves** — `claude-code`, `codex`, `opencode`, `crush`, `antigravity`, `aix` (the `ai-tools` reference pack is pre-installed)
- **Language runtimes the agent will reach for** — Go, Node (via fnm), Bun, Python (via uv), Rust, Java (Temurin LTS), all on PATH with no further setup
- **Everyday CLI building blocks** — bat, fd, ripgrep, lsd, jq, yq, fzf, gh, zoxide, gron, sd, starship, plus tmux + neovim
- **Cloud + security tooling** — aws/azure/gcloud CLIs, kubectl, terraform, trivy, nuclei, httpx, dnsx, subfinder, ffuf, katana, and the rest of the security/cloudsec/appsec packs
- **Sandbox isolation** — everything runs as the non-root `cps` user inside a disposable container; `sudo` is available for ad-hoc package installs without polluting your host

```bash
docker run -d --name agent-sandbox tanq16/cps-sandbox:latest
docker exec -it agent-sandbox zsh -l
# inside the container:
claude-code   # or codex, opencode, crush, ...
```

This is the use case the image is tuned for: an agent (or a human delegating to one) lands in a shell where every tool it's likely to invoke — for code, search, package management, cloud ops, scanning, or recon — is already on PATH. No `brew install` round-trips, no runtime bootstrapping, no "let me set up your environment first." For ephemeral runs, add `--rm` to `docker run`; for sessions you want to come back to, keep the container around and re-`exec` in.

## Deep Removal

Run the included script to wipe CPS, CPS-installed brew packages, and Oh My Zsh. Homebrew itself and `~/.zsh_history` are preserved so you can reinstall cleanly without rebuilding your shell history or re-bootstrapping brew.

```bash
./deep-removal.sh
```

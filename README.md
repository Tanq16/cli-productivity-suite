<div align="center">
  <img src=".github/assets/logo.svg" alt="CLI Productivity Suite Logo" width="200">
  <h1>CLI Productivity Suite</h1>

  <a href="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanq16/cli-productivity-suite/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://github.com/tanq16/cli-productivity-suite/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanq16/cli-productivity-suite"></a><br><br>

  <a href="#capabilities">Capabilities</a> &bull; <a href="#install">Install</a> &bull; <a href="#usage">Usage</a> &bull; <a href="#notes">Notes</a>
</div>

---

A single Go binary (`cps`) that installs and maintains a complete CLI development environment on Linux and macOS: zsh, Neovim, tmux, their configs, and on top of that language runtimes, cloud CLIs, security tooling, and self-hosted services.

It is not a dotfiles manager and not a plugin framework. CPS owns what it installs, all of it under `~/shell/`, and it tracks versions so re-running a command updates rather than reinstalls.

## Capabilities

| Category | Commands | Description |
|---|---|---|
| Setup | `prereq`, `shell` | System prerequisites, then the shell environment: binaries, Neovim, tmux, kitty and configs |
| Packages | `package list`, `package install <name>` | Install a package, a whole group, or a runtime suite |
| Inspection | `status [--check] [--json]` | What is installed, at what version, and what is out of date |
| Maintenance | `system update`, `system kitty`, `self-update` | Upgrade the OS packages, install kitty, replace the `cps` binary |
| Appearance | `theme [name]` | Switch the colour palette across kitty, tmux, Neovim and the CLI tools at once |
| Reference | `cheat <topic>` | Terminal cheat sheets for twelve topics |
| Escape hatch | `custom [args...]` | Run your own `~/.config/cps/custom.sh` with the CPS PATH |

## Install

### Binary

```bash
ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH=amd64; [ "$ARCH" = "aarch64" ] && ARCH=arm64
mkdir -p "$HOME/.local/bin"
curl -sL "https://github.com/tanq16/cli-productivity-suite/releases/latest/download/cps-$(uname -s | tr '[:upper:]' '[:lower:]')-$ARCH" -o "$HOME/.local/bin/cps"
chmod +x "$HOME/.local/bin/cps"
```

Every release carries `linux` and `darwin` builds for `amd64` and `arm64`. On a fresh machine `~/.local/bin` is not on PATH yet, so run the first command by full path as `~/.local/bin/cps prereq`; `cps shell` adds the directory for every session after that.

### From source

Needs Go 1.27.

```bash
git clone https://github.com/tanq16/cli-productivity-suite && cd cli-productivity-suite
make build
```

### Requirements

`cps prereq` installs everything CPS needs, which is why it is the first command you run. It uses `sudo`, and it is the only command besides `cps system` that does.

[Kitty](https://sw.kovidgoyal.net/kitty/) is recommended rather than required, and `cps system kitty` installs it. `cps theme` only writes a kitty palette, and the kitty, tmux and Neovim configs draw glyphs a plain font renders as boxes. On macOS `cps package install nerd-font` provides one; on Linux, install a Nerd Font yourself.

## Usage

Every command is idempotent. Installed versions are tracked in `~/.config/cps/state.json` and anything already current is skipped, so re-running a command is how you update. Binaries land in `~/shell/extensions/`, which `cps shell` puts on PATH.

`--debug` applies to every command and turns the styled output into verbose logging, JSON when piped. `cps package install` and `cps status` also take `--gh-token`, a GitHub PAT for private repos, which defaults to `GITHUB_TOKEN` or `GH_TOKEN` and falls back to `gh auth token` when `gh` is authenticated.

### First run

```bash
cps prereq                      # system packages and Homebrew, uses sudo
cps shell                       # shell binaries, Neovim, plugins, configs
cps package install runtime     # and whatever else you want
```

Run them in that order. `cps shell` deploys a `.zshrc` that reaches for the binaries it installs in the same run, and everything past those two is optional. `scripts/bootstrap.sh` runs the whole sequence unattended for a fresh cloud box or container.

`cps shell` installs 26 packages: 15 CLI binaries, Neovim, the two zsh plugins, and the eight config files. They are not addressable through `cps package`, and updating them means re-running `cps shell`.

### `cps package install <package|group|suite>`

```bash
cps package list                      # every group, suite and package
cps package install binaries          # a whole group
cps package install js-suite          # a runtime suite
cps package install nuclei            # one package
```

| Group | Contents |
|---|---|
| `brew` | Brew packages: ffmpeg, imagemagick, nmap, openssl, and the AWS, Azure and gcloud CLIs |
| `macos-desktop` | aerospace, sol, and the JetBrains Mono Nerd Font. macOS only |
| `runtime` | Go, Node, Python, Java and Rust, each as a suite, with bun and the gopls, pyright, typescript-language-server and ruff language servers |
| `ai` | claude-code, codex, antigravity, cursor-agent |
| `binaries` | nuclei, nuclei-templates, naabu, subfinder, proxify, httpx, dnsx, katana, trufflehog, ffuf, gobuster, gau, gowitness, kubelogin, grpcurl, terraform, kubectl, trivy |
| `homelab` | caddy, senkaimon, linksnapper, kairo, raikiri, expenseowl, backhub, local-content-share, goff, yt-dlp, telly, and the rinnegan, code-server and neo4j app bundles |
| `private` | anbu, box, gcli, nits, sharingan, claudex, toon, cybernest |

A package that needs another is installed after it, and a missing dependency is pulled in rather than failing at the package manager: `cps package install claude-code` installs Node first. The five runtime suites are `go-suite`, `js-suite`, `python-suite`, `java-suite` and `rust-suite`, and each is installable on its own.

`telly`, `toon` and `cybernest` live in private repos and need `--gh-token`. Without one they are skipped with a warning and the rest of the group still installs.

### `cps status`

```bash
cps status                # installed or not, and the recorded version
cps status --check        # compare every installed package against upstream
cps status --json         # the same data, for a script
```

> [!WARNING]
> `cps status --check` makes one GitHub API call per binary package. Without a token this exhausts the 60-per-hour unauthenticated limit in a single run.

Packages whose version is owned by a package manager report it as such: `brew-managed`, `npm-managed`, `uv-managed`, `custom-managed` or `deployed`. Re-running their install is the update path.

### `cps system <verb>`

```bash
cps system update         # apt upgrade and autoremove on Linux, brew upgrade and cleanup on both
cps system kitty          # run the kitty terminal installer
```

`cps system update` is what keeps brew packages current. `HOMEBREW_NO_AUTO_UPDATE=1` is set in the generated rc file, so `brew install` never refreshes its index on its own and a re-run of `cps package install brew` upgrades only against a stale one.

### `cps custom [args...]`

Runs `~/.config/cps/custom.sh` with the CPS PATH in front, passes the arguments through verbatim, and exits with the script's exit code. Everything else about it is yours. [docs/custom-scripts.md](docs/custom-scripts.md) has the contract and an example.

### `cps theme [name]`

```bash
cps theme                 # list, marking the active one
cps theme gruvbox-dark    # switch
```

Fifteen themes ship embedded, in dark and light pairs where upstream publishes both: `mocha` and `latte` (Catppuccin), `gruvbox`, `dracula`, `tokyonight`, `monokai`, `atom-one`, `everforest`, and `nord-dark`.

Only the kitty palette changes, and everything downstream follows because nothing downstream names a hex colour. The tmux, Neovim and starship configs reference ANSI indices `0-15`, and bat, lsd and fzf are pinned to those same indices rather than the 256-colour cube they otherwise default to. Kitty reloads within a second, running tmux and Neovim sessions pick the palette up on their next redraw, and nothing restarts. The choice is recorded in state, so `cps shell` keeps it.

### `cps cheat <topic>`

Cheat sheets for `cps`, `go`, `java`, `uv`, `fnm`, `bun`, `rust`, `tmux`, `nvim`, `fzf`, `jq` and `regex`. `cps cheat list` prints the set with descriptions.

### `cps self-update`

Downloads the latest release and replaces the running binary in place, at whatever path it is running from.

## Notes

- **One generated rc file, not a set of fragments.** `~/.zshrc` sources `~/shell/rc/cps.zsh`, which CPS rewrites whole on every `cps shell` and every `cps package install`. Which snippets it contains follows what is installed, so a runtime or cloud CLI wires itself up the moment it lands. Editing it is pointless because the next run overwrites it. [docs/shell-environment.md](docs/shell-environment.md) covers what it sets.
- **Your own aliases and binaries have two drop-zones.** `.zsh` files in `~/shell/rc/custom/` are sourced after `cps.zsh`, so they can override anything CPS set, and `~/shell/custom-bin/` is prepended to PATH ahead of the CPS-managed directories, so your binary wins a name collision. `cps shell` creates both and never touches them again.
- **Binaries and app bundles land in different places.** Single binaries go to `~/shell/extensions/`, which is on PATH. The multi-file trees (Neovim, rinnegan, code-server, neo4j) unpack to `~/shell/apps/<name>/`, and only Neovim is symlinked onto PATH. [docs/app-bundles.md](docs/app-bundles.md) covers running the others and where each keeps its data.
- **Snapshot directories are replaced wholesale.** The zsh plugins and `nuclei-templates` are downloaded as tarballs of the default branch with no `.git` at all, and an update swaps the whole directory. Keep nothing of your own in `~/shell/plugins/` or `~/shell/nuclei-templates`; custom nuclei templates belong in a directory of your own that you pass with `-t`.
- **Sudo is confined to two commands.** `cps prereq` and `cps system` use it because installing system packages needs it. Nothing in `cps shell` or `cps package` escalates.
- **Removal is a script, not a command.** `./scripts/deep-removal.sh` wipes the `cps` binary, the whole `~/shell/` tree, the CPS-deployed configs, and the brew packages CPS installed. Homebrew itself, `~/.zsh_history`, the macOS desktop casks and `~/.config/neo4j/` are left alone.

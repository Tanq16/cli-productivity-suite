# The shell environment

What `cps shell` puts on disk, how the zsh config is assembled, and where your own additions go.

## The `~/shell/` tree

Everything CPS owns lives here, so removing it removes CPS.

| Directory | Holds |
|---|---|
| `extensions/` | Every single-binary tool, on PATH |
| `custom-bin/` | Your own binaries, on PATH ahead of `extensions/` so yours wins a name collision |
| `apps/` | Multi-file app bundles, not on PATH |
| `rc/` | The generated `cps.zsh`, which is the whole CPS zsh config |
| `rc/custom/` | Your own zsh files, sourced after `cps.zsh` so they can override it |
| `env/` | Generated environment, currently `brew shellenv` output |
| `completions/` | Generated zsh completions for fzf, uv, fnm, zoxide and starship |
| `plugins/` | The two zsh plugins, as snapshots of their default branch |

Language runtimes install here too, rather than into system locations: `go-sdk/`, `java-sdk/`, `rust/`, `fnm/`, `bun/`, `py-default/`, `uv-tools/` and `uv-python/`.

`cps shell` creates `custom-bin/` and `rc/custom/` and nothing CPS runs later touches them. `scripts/deep-removal.sh` is the exception: it wipes the whole `~/shell/` tree, so anything you want to survive a reinstall belongs outside it.

## The config files

`cps shell` deploys these outside `~/shell/`, overwriting whatever is at each path:

| File | Purpose |
|---|---|
| `~/.zshrc` | Thin loader, nothing else |
| `~/.tmux.conf` | The whole tmux configuration |
| `~/.config/nvim/init.lua` | The whole Neovim configuration |
| `~/.config/kitty/kitty.conf` | Kitty settings |
| `~/.config/kitty/current-theme.conf` | The active palette, rewritten by `cps theme` |
| `~/.config/lsd/config.yaml`, `colors.yaml` | lsd settings pinned to ANSI indices |
| `~/.config/starship.toml` | The starship prompt |
| `~/.aerospace.toml` | Aerospace window manager, macOS only |

## One generated file, not a monolithic `.zshrc`

`~/.zshrc` is a loader and nothing more. It sources `~/shell/rc/cps.zsh`, then every `.zsh` file in `~/shell/rc/custom/`.

`cps.zsh` is generated whole on every `cps shell` and every `cps package install`. CPS owns it, so an edit to it is lost on the next run and your own zsh belongs in `~/shell/rc/custom/`. Generation reads the installed set out of `state.json` and appends each snippet whose condition holds, in a fixed order:

| Snippet | Present when | Sets |
|---|---|---|
| base | always | Homebrew, plugins, history, completion, keybindings, PATH, aliases, fzf, zoxide |
| `alias a=anbu` | `anbu` is installed | the `a` alias |
| runtimes | any runtime is installed | `GOROOT`, `GOPATH`, `JAVA_HOME`, `CARGO_HOME`, `FNM_DIR`, `BUN_INSTALL`, `VIRTUAL_ENV`, the `UV_*` paths, and `fnm env` |
| cloud | `aws-cli` or `gcloud-cli` is installed | gcloud path, lazily loaded aws and gcloud completions |
| security | `nuclei-templates` is installed | `NUCLEI_TEMPLATES_DIR` |
| homelab | `neo4j` is installed | `NEO4J_CONF` |
| syntax highlighting | always | sources the zsh-syntax-highlighting plugin |

Syntax highlighting is last because upstream requires it, and generating the file is what makes that ordering reliable.

The aws and gcloud completions are wrapped in a shim that loads them on first invocation, because sourcing the gcloud one eagerly costs about 30ms of a roughly 75ms startup.

## Neovim

Neovim installs as an app bundle under `~/shell/apps/neovim/`, with `~/shell/extensions/nvim` symlinked to its launcher, because Neovim locates its own runtime relative to the resolved executable.

The config is a single `init.lua` leaning on Neovim natives rather than plugins: `vim.pack` for plugin management, the built-in LSP client, autotriggered built-in completion on any server that supports it, and `listchars` indent guides. Five plugins fill the gaps: `fzf-lua`, `nvim-tree`, `gitsigns`, `nvim-autopairs`, and `nvim-treesitter`.

Plugins and parsers install on first launch, in the background. Treesitter's `main` branch shells out to the `tree-sitter` CLI to compile each parser, so parsers for bash, go, gomod, javascript, json, lua, markdown, python, tsx, typescript and yaml only appear once `cps shell` has provided that CLI. Without it, buffers fall back to the bundled regex syntax rather than failing.

Four language servers are configured, each enabled only when its binary is on PATH, and `cps package install runtime` installs all four:

| Server | Covers |
|---|---|
| `gopls` | Go, gomod, gowork, gotmpl |
| `pyright` | Python types, hover, go-to-definition |
| `ruff` | Python linting and formatting |
| `ts_ls` | JavaScript and TypeScript, including the React filetypes |

Python gets two because they cover different ground and neither does the other's job.

Colours come from the terminal. `termguicolors` is off, so highlights resolve against ANSI 0-15 and follow whatever `cps theme` set. The statusline is hand-rolled from those same indices, showing mode, git branch, path, diagnostics, filetype and position, and capped with the glyphs the tmux bar uses so the two read as one design.

## tmux

No plugins and no plugin manager. `.tmux.conf` is the whole configuration.

The status bar sits at the top and uses ANSI indices with `bg=default` throughout, so it stays transparent against the terminal background and re-themes with it. The right side shows the current directory and the session name, and the session glyph turns red while the prefix is pending. `default-terminal` is pinned to `tmux-256color` so italics and undercurl reach Neovim intact. The glyphs need a Nerd Font.

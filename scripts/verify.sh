#!/usr/bin/env bash
# verify.sh - sanity check for a fully installed cps environment.
#
# Verifies rc loading, directories, files, env vars, PATH segments, and the
# binary each group installs. Silent on success; prints only failures.
# Exits 0 if all good, 1 otherwise.
#
# Usage:
#   zsh -lc 'bash scripts/verify.sh'
#
# The zsh -l wrapper is required: it sources the login shell, which loads
# ~/shell/rc/cps.zsh and exports the env vars this script checks.
#
# It assumes every group is installed. Skip a group's section by hand when
# checking a partial install.

failed=()
fail() { failed+=("$1"); }
check_bin() { command -v "$1" >/dev/null 2>&1 || fail "$2: $1"; }

case "$(uname -s)" in
    Darwin) OS=darwin ;;
    *) OS=linux ;;
esac

# --- directories ---
for d in \
    shell/extensions shell/custom-bin shell/plugins shell/apps \
    shell/rc shell/rc/custom shell/env shell/completions \
    shell/go-sdk shell/java-sdk shell/rust shell/fnm shell/py-default \
    shell/uv-tools shell/uv-tool-executables shell/uv-python \
    shell/nuclei-templates \
    .config/cps .config/nvim; do
    [ -d "$HOME/$d" ] || fail "dir missing: ~/$d"
done

# --- files ---
for f in \
    .zshrc .tmux.conf \
    .config/kitty/kitty.conf .config/kitty/current-theme.conf \
    .config/starship.toml .config/cps/state.json \
    .config/lsd/colors.yaml .config/lsd/config.yaml \
    .config/nvim/init.lua \
    shell/apps/neovim/bin/nvim \
    shell/rc/cps.zsh \
    shell/env/brew.zsh \
    shell/completions/fzf.zsh shell/completions/uv.zsh \
    shell/completions/fnm.zsh shell/completions/zoxide.zsh \
    shell/completions/starship.zsh; do
    [ -f "$HOME/$f" ] || fail "file missing: ~/$f"
done

for stale in 00-base.zsh 10-runtimes.zsh 20-cloud.zsh 30-security.zsh 40-misc.zsh 50-homelab.zsh; do
    [ -f "$HOME/shell/rc/$stale" ] && fail "stale rc fragment: ~/shell/rc/$stale"
done

# brew.zsh must be non-empty - guards against the cold-shellenv early-return bug.
[ -s "$HOME/shell/env/brew.zsh" ] || fail "~/shell/env/brew.zsh is empty (cold brew shellenv)"

[ "$(grep -c zsh-syntax-highlighting "$HOME/shell/rc/cps.zsh")" = "1" ] || \
    fail "cps.zsh does not source zsh-syntax-highlighting exactly once"
tail -1 "$HOME/shell/rc/cps.zsh" | grep -q zsh-syntax-highlighting || \
    fail "cps.zsh does not source zsh-syntax-highlighting last"

# --- env vars ---
for v in GOROOT GOPATH JAVA_HOME RUSTUP_HOME CARGO_HOME FNM_DIR \
         FNM_MULTISHELL_PATH BUN_INSTALL npm_config_cache VIRTUAL_ENV \
         UV_TOOL_DIR UV_TOOL_BIN_DIR UV_PYTHON_INSTALL_DIR \
         NUCLEI_TEMPLATES_DIR NEO4J_CONF; do
    eval "val=\$$v"
    [ -n "$val" ] || fail "env unset: $v"
done

# --- PATH composition ---
for p in \
    "$HOME/shell/custom-bin" "$HOME/shell/extensions" \
    "$HOME/shell/uv-tool-executables" "$HOME/shell/go-sdk/bin" \
    "$HOME/shell/go/bin" "$HOME/shell/java-sdk/bin" \
    "$HOME/shell/rust/.cargo/bin" "$HOME/shell/bun/bin" \
    "$HOME/shell/py-default/bin"; do
    case ":$PATH:" in
        *":$p:"*) ;;
        *) fail "PATH missing: $p" ;;
    esac
done

# --- shell ---
for t in bat fd rg lsd jq yq fzf gh tree-sitter gron zoxide sd starship age sq; do
    check_bin "$t" "shell"
done
# nvim reaches PATH via a symlink; the bundle itself must stay intact
check_bin nvim "shell"
[ "$(readlink "$HOME/shell/extensions/nvim")" = "$HOME/shell/apps/neovim/bin/nvim" ] || \
    fail "shell: ~/shell/extensions/nvim does not point at the neovim bundle"
for p in zsh-autosuggestions zsh-syntax-highlighting; do
    [ -f "$HOME/shell/plugins/$p/$p.zsh" ] || fail "shell: ~/shell/plugins/$p/$p.zsh"
done

# --- brew ---
for t in ffmpeg magick nmap aws az gcloud; do check_bin "$t" "brew"; done

# --- macos-desktop ---
if [ "$OS" = "darwin" ]; then
    check_bin aerospace "macos-desktop"
    [ -d "/Applications/Sol.app" ] || fail "macos-desktop: /Applications/Sol.app"
    ls "$HOME/Library/Fonts/JetBrainsMono"*NerdFont* >/dev/null 2>&1 || \
        fail "macos-desktop: no JetBrains Mono Nerd Font in ~/Library/Fonts"
fi

# --- runtime ---
for t in uv fnm bun go java python rustc cargo node npm \
         gopls pyright typescript-language-server ruff; do
    check_bin "$t" "runtime"
done

# --- ai ---
for t in claude codex cursor-agent agy; do check_bin "$t" "ai"; done

# --- binaries ---
for t in nuclei naabu subfinder proxify httpx dnsx katana trufflehog \
         ffuf gobuster gau gowitness \
         kubelogin grpcurl terraform kubectl trivy; do
    check_bin "$t" "binaries"
done
[ -n "$(ls -A "$HOME/shell/nuclei-templates" 2>/dev/null)" ] || \
    fail "binaries: ~/shell/nuclei-templates is empty"

# --- homelab ---
for t in caddy senkaimon linksnapper kairo raikiri expenseowl backhub \
         local-content-share goff yt-dlp telly; do
    check_bin "$t" "homelab"
done
# app bundles are not on PATH, so check the launcher inside the bundle
for b in rinnegan/bin/rinnegan code-server/bin/code-server \
         neo4j/bin/neo4j neo4j/bin/cypher-shell; do
    [ -x "$HOME/shell/apps/$b" ] || fail "homelab: ~/shell/apps/$b"
done
# neo4j needs a JVM, and its state must live outside the bundle or an upgrade
# destroys the databases
[ -x "$JAVA_HOME/bin/java" ] || fail "homelab: neo4j has no JVM at \$JAVA_HOME/bin/java"
[ -f "$NEO4J_CONF/neo4j.conf" ] || fail "homelab: no neo4j.conf at \$NEO4J_CONF"
for d in data plugins import logs run licenses; do
    [ -d "$HOME/.config/neo4j/$d" ] || fail "homelab: ~/.config/neo4j/$d missing"
    grep -q "^server.directories.$d=$HOME/.config/neo4j/$d\$" "$NEO4J_CONF/neo4j.conf" 2>/dev/null || \
        fail "homelab: neo4j.conf does not relocate $d"
done
# code-server must be pre-seeded; without it the first launch generates a password config
[ -f "$HOME/.config/code-server/config.yaml" ] || fail "homelab: ~/.config/code-server/config.yaml missing"
grep -q '^auth: none$' "$HOME/.config/code-server/config.yaml" 2>/dev/null || \
    fail "homelab: code-server auth is not disabled"
[ -f "$HOME/.local/share/code-server/User/settings.json" ] || \
    fail "homelab: code-server settings.json missing"

# --- private ---
for t in anbu box gcli nits sharingan claudex toon cybernest; do
    check_bin "$t" "private"
done

# --- cps agrees ---
command -v cps >/dev/null 2>&1 && {
    missing=$(cps status --json | jq -r '.[] | select(.installed == false) | .name')
    [ -z "$missing" ] || fail "cps status reports not installed: $(echo "$missing" | tr '\n' ' ')"
}

# --- result ---
if [ ${#failed[@]} -eq 0 ]; then
    exit 0
fi

printf '%s\n' "${failed[@]}"
echo ""
echo "FAILURES: ${#failed[@]}"
exit 1

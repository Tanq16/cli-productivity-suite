#!/usr/bin/env bash
# verify.sh - sanity check for a built cps environment.
#
# Verifies identity, rc loading, directories, files, env vars, PATH segments,
# and binaries across every pack. Silent on success; prints only failures.
# Exits 0 if all good, 1 otherwise.
#
# Usage:
#   zsh -lc 'bash scripts/verify.sh'
#
# The zsh -l wrapper is required: it sources ~/.zprofile, which loads the cps
# rc fragments and exports the env vars this script checks.

failed=()
fail() { failed+=("$1"); }

# --- identity ---
[ "$(whoami)" = "cps" ] || fail "user: expected 'cps', got '$(whoami)'"
[ "$(id -u)" = "1000" ] || fail "uid: expected 1000, got $(id -u)"
sudo -n true 2>/dev/null || fail "sudo: NOPASSWD not configured"

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
    .zshrc .zprofile .tmux.conf \
    .config/kitty/kitty.conf .config/kitty/current-theme.conf \
    .config/starship.toml .config/cps/state.json \
    .config/lsd/colors.yaml .config/lsd/config.yaml \
    .config/nvim/init.lua \
    shell/apps/neovim/bin/nvim \
    shell/rc/00-base.zsh shell/rc/10-runtimes.zsh \
    shell/rc/20-cloud.zsh shell/rc/30-security.zsh shell/rc/50-homelab.zsh \
    shell/env/brew.zsh \
    shell/completions/fzf.zsh shell/completions/uv.zsh \
    shell/completions/fnm.zsh shell/completions/zoxide.zsh \
    shell/completions/starship.zsh; do
    [ -f "$HOME/$f" ] || fail "file missing: ~/$f"
done

# brew.zsh must be non-empty - guards against the cold-shellenv early-return bug.
[ -s "$HOME/shell/env/brew.zsh" ] || fail "~/shell/env/brew.zsh is empty (cold brew shellenv)"

# --- env vars ---
for v in GOROOT GOPATH JAVA_HOME RUSTUP_HOME CARGO_HOME FNM_DIR \
         FNM_MULTISHELL_PATH BUN_INSTALL npm_config_cache VIRTUAL_ENV \
         UV_TOOL_DIR UV_TOOL_BIN_DIR UV_PYTHON_INSTALL_DIR; do
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

# --- binaries by pack ---
check_bin() { command -v "$1" >/dev/null 2>&1 || fail "$2: $1"; }

# essentials
for t in bat fd rg lsd jq yq fzf gh gron zoxide sd starship anbu danzo; do
    check_bin "$t" "essentials"
done
# base app bundle - nvim reaches PATH via a symlink; the bundle itself must stay intact
check_bin nvim "base"
[ "$(readlink "$HOME/shell/extensions/nvim")" = "$HOME/shell/apps/neovim/bin/nvim" ] || \
    fail "base: ~/shell/extensions/nvim does not point at the neovim bundle"
# runtimes
for t in uv fnm bun go java python rustc cargo node npm; do
    check_bin "$t" "runtimes"
done
# cloud
for t in aws az gcloud; do check_bin "$t" "cloud"; done
# security
for t in nuclei naabu subfinder proxify trufflehog httpx dnsx; do
    check_bin "$t" "security"
done
# cloudsec
for t in kubelogin grpcurl terraform kubectl trivy prowler oci tofu; do
    check_bin "$t" "cloudsec"
done
# appsec tools now ship in security
for t in katana ffuf dalfox gobuster gau; do
    check_bin "$t" "security"
done
# gowitness moved to security; age and sq to essentials
check_bin gowitness "security"
for t in age sq; do check_bin "$t" "essentials"; done
# homelab app bundle - neo4j needs a JVM, so check the launcher and that JAVA_HOME resolves one
[ -x "$HOME/shell/apps/neo4j/bin/neo4j" ] || fail "homelab: ~/shell/apps/neo4j/bin/neo4j"
[ -x "$HOME/shell/apps/neo4j/bin/cypher-shell" ] || fail "homelab: ~/shell/apps/neo4j/bin/cypher-shell"
[ -x "$JAVA_HOME/bin/java" ] || fail "homelab: neo4j has no JVM at \$JAVA_HOME/bin/java"
# neo4j state must live outside the bundle, or an upgrade destroys the databases
[ -n "$NEO4J_CONF" ] || fail "homelab: NEO4J_CONF unset - neo4j would use the bundle's own conf"
[ -f "$NEO4J_CONF/neo4j.conf" ] || fail "homelab: no neo4j.conf at \$NEO4J_CONF"
for d in data plugins import logs run licenses; do
    [ -d "$HOME/.config/neo4j/$d" ] || fail "homelab: ~/.config/neo4j/$d missing"
    grep -q "^server.directories.$d=$HOME/.config/neo4j/$d\$" "$NEO4J_CONF/neo4j.conf" 2>/dev/null || \
        fail "homelab: neo4j.conf does not relocate $d"
done
# ai-tools
for t in claude codex cursor-agent agy; do check_bin "$t" "ai-tools"; done
# homelab
for t in caddy linksnapper kairo raikiri expenseowl backhub local-content-share goff yt-dlp; do check_bin "$t" "homelab"; done
# homelab app bundles - not on PATH, so check the launcher inside the bundle
for b in rinnegan/bin/rinnegan code-server/bin/code-server; do
    [ -x "$HOME/shell/apps/$b" ] || fail "homelab: ~/shell/apps/$b"
done
# code-server must be pre-seeded; without it the first launch generates a password config
[ -f "$HOME/.config/code-server/config.yaml" ] || fail "homelab: ~/.config/code-server/config.yaml missing"
grep -q '^auth: none$' "$HOME/.config/code-server/config.yaml" 2>/dev/null || \
    fail "homelab: code-server auth is not disabled"
[ -f "$HOME/.local/share/code-server/User/settings.json" ] || \
    fail "homelab: code-server settings.json missing"
# private (public subset)
for t in nits gcli box claudex sharingan; do check_bin "$t" "private"; done

# --- nuclei templates non-empty ---
[ -n "$(ls -A "$HOME/shell/nuclei-templates" 2>/dev/null)" ] || \
    fail "~/shell/nuclei-templates is empty"

# --- result ---
if [ ${#failed[@]} -eq 0 ]; then
    exit 0
fi

printf '%s\n' "${failed[@]}"
echo ""
echo "FAILURES: ${#failed[@]}"
exit 1

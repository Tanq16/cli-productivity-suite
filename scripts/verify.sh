#!/usr/bin/env bash
# verify.sh — sanity check for a built cps environment.
#
# Verifies identity, rc loading, directories, files, env vars, PATH segments,
# and binaries across every pack. Silent on success; prints only failures.
# Exits 0 if all good, 1 otherwise.
#
# Usage (against a running container):
#   docker exec -it <container> zsh -lc 'bash /path/to/verify.sh'
#
# Usage (one-shot against an image):
#   docker run --rm -v "$PWD/scripts/verify.sh:/tmp/v.sh" \
#       tanq16/cps-sandbox:latest zsh -lc 'bash /tmp/v.sh'
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
    shell/executables shell/extensions shell/custom-bin shell/plugins shell/apps \
    shell/rc shell/rc/custom shell/env shell/completions \
    shell/go-sdk shell/java-sdk shell/rust shell/fnm shell/py-default \
    shell/uv-tools shell/uv-tool-executables shell/uv-python \
    shell/nuclei-templates \
    .config/cps .config/nvim .tmux/plugins/tpm; do
    [ -d "$HOME/$d" ] || fail "dir missing: ~/$d"
done

# --- files ---
for f in \
    .zshrc .zprofile .tmux.conf \
    .config/kitty/kitty.conf .config/kitty/current-theme.conf \
    .config/starship.toml .config/cps/state.json \
    shell/rc/00-base.zsh shell/rc/10-runtimes.zsh \
    shell/rc/20-cloud.zsh shell/rc/30-security.zsh \
    shell/env/brew.zsh \
    shell/completions/fzf.zsh shell/completions/uv.zsh \
    shell/completions/fnm.zsh shell/completions/zoxide.zsh \
    shell/completions/starship.zsh; do
    [ -f "$HOME/$f" ] || fail "file missing: ~/$f"
done

# brew.zsh must be non-empty — guards against the cold-shellenv early-return bug.
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
    "$HOME/shell/custom-bin" "$HOME/shell/extensions" "$HOME/shell/executables" \
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
# appsec
for t in katana ffuf dalfox gobuster gau; do
    check_bin "$t" "appsec"
done
# misc
for t in gowitness age sq; do check_bin "$t" "misc"; done
# ai-tools
for t in claude codex cursor-agent agy; do check_bin "$t" "ai-tools"; done
# homelab
for t in caddy linksnapper kairo raikiri expenseowl; do check_bin "$t" "homelab"; done
# homelab app bundles — not on PATH, so check the launcher inside the bundle
for b in rinnegan/bin/rinnegan code-server/bin/code-server; do
    [ -x "$HOME/shell/apps/$b" ] || fail "homelab: ~/shell/apps/$b"
done
# private (public subset)
for t in nits gcli box claudex; do check_bin "$t" "private"; done

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

#!/usr/bin/env bash

EXT_BINS=(
  ai-context aix aurelian cloudfox cloudlist hadrian julius lincli
  nblm nerva poltergeist reaper snitch titus trajan wraith
)
GOBIN_BINS=(augustus vespasian)
UV_TOOLS=(checkov mycli pgcli praetorian-cli)
NPM_PKGS=(opencode-ai @charmland/crush)
BREW_PKGS=(usql)
STATE_KEYS=(
  ai-context aix augustus aurelian checkov cloudfox cloudlist crush
  hadrian julius lincli mycli nblm nerva opencode pgcli poltergeist
  praetorian-cli reaper snitch titus trajan usql vespasian wraith
)

echo "CPS 1.8.0 migration"
echo ""
echo "Will remove:"
echo "  - ~/shell/executables (no longer created; only if empty)"
echo "  - ~/.config/cps/extensions (custom-extension YAML cache, system removed)"
echo "  - Dropped binaries from ~/shell/extensions and ~/shell/go/bin"
echo "  - Dropped uv tools (${UV_TOOLS[*]})"
echo "  - Dropped npm globals (${NPM_PKGS[*]})"
echo "  - Dropped brew formulas (${BREW_PKGS[*]})"
echo "  - Their entries in ~/.config/cps/state.json"
echo ""
echo "Will preserve:"
echo "  - Everything still shipped in 1.8.0 (prowler, oci-cli, tofu, sq, caddy, ...)"
echo "  - ~/shell/custom-bin and ~/shell/rc/custom (your drop zones)"
echo "  - Anything you placed in ~/shell/executables yourself"
echo ""
read -rp "Continue? [y/N] " ans
case "$ans" in
  [yY]|[yY][eE][sS]) ;;
  *) echo "aborted"; exit 0 ;;
esac

echo ""
echo "==> removing dropped binaries from ~/shell/extensions"
for b in "${EXT_BINS[@]}"; do
  rm -f "$HOME/shell/extensions/$b"
done

echo "==> removing dropped binaries from ~/shell/go/bin"
for b in "${GOBIN_BINS[@]}"; do
  rm -f "$HOME/shell/go/bin/$b"
done

if command -v uv >/dev/null 2>&1; then
  echo "==> uninstalling dropped uv tools"
  for t in "${UV_TOOLS[@]}"; do
    uv tool uninstall "$t" >/dev/null 2>&1 || true
  done
else
  echo "==> uv not found, skipping uv tool uninstall"
fi

if command -v npm >/dev/null 2>&1; then
  echo "==> uninstalling dropped npm globals"
  for p in "${NPM_PKGS[@]}"; do
    npm uninstall -g "$p" >/dev/null 2>&1 || true
  done
else
  echo "==> npm not found, skipping npm uninstall"
fi

if command -v brew >/dev/null 2>&1; then
  echo "==> uninstalling dropped brew formulas"
  brew uninstall "${BREW_PKGS[@]}" >/dev/null 2>&1 || true
else
  echo "==> brew not found, skipping brew uninstall"
fi

echo "==> removing custom-extension YAML cache"
rm -rf "$HOME/.config/cps/extensions"

echo "==> removing ~/shell/executables"
if [ -d "$HOME/shell/executables" ]; then
  if [ -z "$(ls -A "$HOME/shell/executables" 2>/dev/null)" ]; then
    rmdir "$HOME/shell/executables"
  else
    echo "    NOT removed: ~/shell/executables is not empty"
    echo "    it is no longer on PATH; move what you need to ~/shell/custom-bin"
  fi
fi

STATE="$HOME/.config/cps/state.json"
if [ -f "$STATE" ]; then
  if command -v jq >/dev/null 2>&1; then
    echo "==> pruning state.json"
    keys=$(printf '%s\n' "${STATE_KEYS[@]}" | jq -R . | jq -s .)
    if tmp=$(mktemp) && jq --argjson rm "$keys" \
        'if .tools then .tools |= with_entries(select(.key as $k | $rm | index($k) | not)) else . end' \
        "$STATE" > "$tmp" 2>/dev/null && [ -s "$tmp" ]; then
      mv "$tmp" "$STATE"
    else
      rm -f "$tmp"
      echo "    skipped: could not rewrite state.json"
    fi
  else
    echo "==> jq not found, skipping state.json prune"
  fi
fi

echo ""
echo "done."
echo ""
echo "next:"
echo "  1. open a new shell (or re-source ~/.zshrc) to pick up the new PATH"
echo "  2. cps extend list"

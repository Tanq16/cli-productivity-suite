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
if grep -q 'shell/executables' "$HOME/shell/rc/00-base.zsh" 2>/dev/null; then
  echo "WARNING: ~/shell/rc/00-base.zsh still puts ~/shell/executables on PATH."
  echo "Run 'cps init' with the 1.8.0 binary first, then re-run this script."
  echo ""
fi
echo "Will remove:"
echo "  - ~/shell/executables (offers to move any contents to ~/shell/custom-bin)"
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
echo "  - Existing ~/shell/custom-bin files (never overwritten by the move)"
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

EXECDIR="$HOME/shell/executables"
CUSTOMBIN="$HOME/shell/custom-bin"

echo "==> removing ~/shell/executables"
if [ ! -d "$EXECDIR" ]; then
  echo "    already gone"
elif [ -z "$(ls -A "$EXECDIR" 2>/dev/null)" ]; then
  rmdir "$EXECDIR" && echo "    removed (was empty)"
else
  echo "    not empty, and no longer on PATH:"
  ls -A "$EXECDIR" | while read -r n; do echo "      $n"; done
  read -rp "    move these to ~/shell/custom-bin? [y/N] " mv_ans
  case "$mv_ans" in
    [yY]|[yY][eE][sS])
      mkdir -p "$CUSTOMBIN"
      moved=0
      kept=0
      for f in "$EXECDIR"/* "$EXECDIR"/.[!.]*; do
        [ -e "$f" ] || continue
        base=$(basename "$f")
        if [ -e "$CUSTOMBIN/$base" ]; then
          echo "      kept $base (already in custom-bin)"
          kept=$((kept + 1))
        elif mv "$f" "$CUSTOMBIN/$base"; then
          moved=$((moved + 1))
        else
          kept=$((kept + 1))
        fi
      done
      echo "      moved $moved, kept $kept"
      if [ -z "$(ls -A "$EXECDIR" 2>/dev/null)" ]; then
        rmdir "$EXECDIR" && echo "    removed"
      else
        echo "    kept ~/shell/executables (still has files)"
      fi
      ;;
    *)
      echo "    kept ~/shell/executables; move what you need to ~/shell/custom-bin"
      ;;
  esac
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

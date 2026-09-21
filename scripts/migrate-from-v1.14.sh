#!/usr/bin/env bash
# migrate-from-v1.14.sh - moves yt-dlp off the GitHub release binary and onto a uv tool install.
#
# Through v1.14.x, yt-dlp installed as a PyInstaller bundle in ~/shell/extensions/,
# which unpacks a whole Python environment into a temp directory on every call.
# From the next release it is a uv tool in ~/shell/uv-tool-executables/.
#
# The old binary has to go, not just be left behind: ~/shell/extensions/ sits ahead
# of ~/shell/uv-tool-executables/ on the PATH that cps builds for its subprocesses,
# so a leftover copy keeps winning.
#
# Run this once, after updating cps. Safe to re-run and safe on a machine that
# never had yt-dlp.
#
# Usage:
#   bash scripts/migrate-from-v1.14.sh

set -euo pipefail

STATE="$HOME/.config/cps/state.json"
OLD_BINARY="$HOME/shell/extensions/yt-dlp"

command -v cps >/dev/null 2>&1 || { echo "cps is not on PATH; update it first" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is not on PATH; run \`cps shell\` first" >&2; exit 1; }

tracked=""
if [ -f "$STATE" ]; then
    tracked=$(jq -r '.tools["yt-dlp"].version // ""' "$STATE")
fi

if [ -z "$tracked" ] && [ ! -e "$OLD_BINARY" ]; then
    echo "yt-dlp was never installed here, nothing to migrate"
    exit 0
fi

if [ "$tracked" = "uv-managed" ] && [ ! -e "$OLD_BINARY" ]; then
    echo "yt-dlp is already a uv tool, nothing to migrate"
    exit 0
fi

echo "==> removing the release binary"
rm -f "$OLD_BINARY"

echo "==> reinstalling yt-dlp as a uv tool"
cps package install yt-dlp

echo ""
echo "cps records yt-dlp as $(jq -r '.tools["yt-dlp"].version // "untracked"' "$STATE")"
echo "the executable is now $HOME/shell/uv-tool-executables/yt-dlp"

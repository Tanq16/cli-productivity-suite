#!/usr/bin/env bash
# migrate-from-v1.14.sh - moves yt-dlp onto a uv tool install and gives every uv
# tool a real recorded version.
#
# Through v1.14.x, yt-dlp installed as a PyInstaller bundle in ~/shell/extensions/,
# which unpacks a whole Python environment into a temp directory on every call.
# From the next release it is a uv tool in ~/shell/uv-tool-executables/.
#
# The old binary has to go, not just be left behind: ~/shell/extensions/ sits ahead
# of ~/shell/uv-tool-executables/ on the PATH that cps builds for its subprocesses,
# so a leftover copy keeps winning.
#
# uv tools also recorded the literal string "uv-managed" as their version through
# v1.14.x, so `cps status --check` had nothing to compare against PyPI. Reinstalling
# them writes the real version.
#
# Run this once, after updating cps. Safe to re-run and safe on a machine that
# never had yt-dlp.
#
# Usage:
#   bash scripts/migrate-from-v1.14.sh

set -euo pipefail

OLD_BINARY="$HOME/shell/extensions/yt-dlp"

command -v cps >/dev/null 2>&1 || { echo "cps is not on PATH; update it first" >&2; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is not on PATH; run \`cps shell\` first" >&2; exit 1; }

status=$(cps status --json)
untracked=$(printf '%s' "$status" | jq -r '.[] | select(.kind == "python-tool" and .version == "uv-managed") | .name')

if [ -z "$untracked" ] && [ ! -e "$OLD_BINARY" ]; then
    echo "nothing to migrate"
    exit 0
fi

if [ -e "$OLD_BINARY" ]; then
    echo "==> removing the yt-dlp release binary"
    rm -f "$OLD_BINARY"
fi

installed=$(printf '%s' "$status" | jq -r '.[] | select(.kind == "python-tool" and .installed) | .name')
for tool in $installed; do
    echo "==> reinstalling $tool as a uv tool"
    cps package install "$tool"
done

echo ""
echo "uv tools now record:"
cps status --json | jq -r '.[] | select(.kind == "python-tool" and .installed) | "  \(.name) \(.version)"'

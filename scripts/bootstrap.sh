#!/usr/bin/env bash
# bootstrap.sh - unattended first-run setup for cloud-init and container images.
#
# Runs as the target user, not root. The user needs passwordless sudo, because
# `cps prereq` installs system packages.
#
# Usage:
#   bash scripts/bootstrap.sh
#
# Set GITHUB_TOKEN or GH_TOKEN beforehand to reach the private repos. Without
# one, the private packages are skipped with a warning and the run still
# succeeds.

set -euo pipefail

cps prereq
cps shell

cat > ~/.zprofile <<'EOF'
[[ -o interactive ]] && return
source "$HOME/shell/rc/cps.zsh"
for f in "$HOME/shell/rc/custom/"*.zsh(N); do source "$f"; done
EOF

# Each install runs in a login shell so it inherits the runtime env vars the
# generated rc file exports.
run() { zsh -l -c "$1"; }

run 'cps package install brew'
run 'cps package install runtime'
run 'cps package install binaries'
run 'cps package install homelab'
run 'cps package install private'
run 'cps package install ai'

run 'cps status'

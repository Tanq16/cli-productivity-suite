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

# Each install runs in a login shell, whose ~/.zprofile (set up by cps shell)
# sources the runtime env vars and PATH from cps-env.zsh.
run() { zsh -l -c "$1"; }

run 'cps package install brew'
run 'cps package install runtime'
run 'cps package install binaries'
run 'cps package install homelab'
run 'cps package install private'
run 'cps package install ai'

run 'cps status'

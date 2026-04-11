#!/usr/bin/env bash
# Migration script for CPS v1.3.0
#
# Includes v1.2.0 migrations (tools moved from core to extensions)
# and v1.3.0 changes (consolidating everything under ~/shell/).
#
# v1.2.0: Extension tools moved from ~/shell/executables/ to ~/shell/extensions/
# v1.3.0:
#   - Go SDK:            /usr/local/go        -> ~/shell/go-sdk/
#   - Go workspace:      ~/go/                -> ~/shell/go/
#   - gcloud CLI:        ~/google-cloud-sdk/  -> ~/shell/gcloud-sdk/
#   - nuclei templates:  ~/nuclei-templates/  -> ~/shell/nuclei-templates/
#   - nvm (removed):     ~/.nvm/              -> replaced by fnm in ~/shell/fnm/
#
# Run this script, then run "cps init" to reinstall everything
# in the new locations.

set -euo pipefail

EXEC_DIR="$HOME/shell/executables"
STATE_FILE="$HOME/.config/cps/state.json"

# --- v1.2.0: Tools moved from core to extensions ---
MOVED_TOOLS=(
  # security
  nuclei naabu subfinder proxify trufflehog gobuster
  # appsec
  katana ffuf
  # cloud
  kubelogin grpcurl terraform kubectl
  # private
  nits raikiri gcli box claudex toon cybernest
)

echo "CPS v1.3.0 migration"
echo ""
echo "--- Phase 1: Remove old extension binaries from ~/shell/executables/ ---"

removed=0
for tool in "${MOVED_TOOLS[@]}"; do
  if [ -f "$EXEC_DIR/$tool" ]; then
    rm -f "$EXEC_DIR/$tool"
    echo "  removed binary: $tool"
    ((++removed))
  fi
done

if [ "$removed" -eq 0 ]; then
  echo "  no binaries to remove"
fi

# --- v1.3.0: Remove old directories ---
echo ""
echo "--- Phase 2: Remove old tool locations ---"

removed=0

# nvm -> replaced by fnm
if [ -d "$HOME/.nvm" ]; then
  rm -rf "$HOME/.nvm"
  echo "  removed ~/.nvm/ (replaced by fnm)"
  ((++removed))
fi

# gcloud -> ~/shell/gcloud-sdk/
if [ -d "$HOME/google-cloud-sdk" ]; then
  rm -rf "$HOME/google-cloud-sdk"
  echo "  removed ~/google-cloud-sdk/ (moves to ~/shell/gcloud-sdk/)"
  ((++removed))
fi

# nuclei-templates -> ~/shell/nuclei-templates/
if [ -d "$HOME/nuclei-templates" ]; then
  rm -rf "$HOME/nuclei-templates"
  echo "  removed ~/nuclei-templates/ (moves to ~/shell/nuclei-templates/)"
  ((++removed))
fi

# Go SDK -> ~/shell/go-sdk/
if [ -d "/usr/local/go" ]; then
  sudo rm -rf /usr/local/go
  echo "  removed /usr/local/go (moves to ~/shell/go-sdk/)"
  ((++removed))
fi

# Go workspace -> ~/shell/go/
if [ -d "$HOME/go" ]; then
  mkdir -p "$HOME/shell/go"
  # Preserve module cache and installed binaries
  if [ -d "$HOME/go/bin" ] && [ "$(ls -A "$HOME/go/bin" 2>/dev/null)" ]; then
    cp -r "$HOME/go/bin" "$HOME/shell/go/"
    echo "  migrated ~/go/bin/ -> ~/shell/go/bin/"
  fi
  if [ -d "$HOME/go/pkg" ]; then
    cp -r "$HOME/go/pkg" "$HOME/shell/go/"
    echo "  migrated ~/go/pkg/ -> ~/shell/go/pkg/"
  fi
  rm -rf "$HOME/go"
  echo "  removed ~/go/ (GOPATH moves to ~/shell/go/)"
  ((++removed))
fi

if [ "$removed" -eq 0 ]; then
  echo "  no directories to remove"
fi

# --- Clean state entries ---
echo ""
echo "--- Phase 3: Clean state.json ---"

# Combined: extension tools + relocated tools
ALL_CLEAN_TOOLS=("${MOVED_TOOLS[@]}" nvm go-sdk gcloud-cli nuclei-templates)

if [ -f "$STATE_FILE" ]; then
  if command -v jq &>/dev/null; then
    cleaned=0
    for tool in "${ALL_CLEAN_TOOLS[@]}"; do
      if jq -e ".tools.\"$tool\"" "$STATE_FILE" &>/dev/null; then
        ((++cleaned))
      fi
    done

    if [ "$cleaned" -gt 0 ]; then
      jq_expr="del("
      first=true
      for tool in "${ALL_CLEAN_TOOLS[@]}"; do
        if [ "$first" = true ]; then
          jq_expr+=".tools.\"$tool\""
          first=false
        else
          jq_expr+=", .tools.\"$tool\""
        fi
      done
      jq_expr+=")"

      jq "$jq_expr" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
      echo "  cleaned $cleaned entries from state.json"
    else
      echo "  no state entries to clean"
    fi
  else
    echo "  warning: jq not found, skipping state.json cleanup"
    echo "  install jq and re-run, or manually remove tool entries from $STATE_FILE"
  fi
else
  echo "  no state.json found, skipping"
fi

echo ""
echo "done! now run:"
echo "  cps init"

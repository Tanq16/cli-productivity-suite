#!/usr/bin/env bash
# Migration script for CPS v1.2.0
#
# These tools moved from core (~/shell/executables/) to extensions
# (~/shell/extensions/). This script removes the old binaries and
# cleans their entries from state.json so they can be reinstalled
# cleanly via "cps extend <pack>".

set -euo pipefail

EXEC_DIR="$HOME/shell/executables"
STATE_FILE="$HOME/.config/cps/state.json"

# Tools moved to extensions
MOVED_TOOLS=(
  # security
  nuclei naabu subfinder proxify trufflehog
  # appsec
  katana ffuf
  # cloud
  kubelogin grpcurl terraform kubectl
  # private
  nits raikiri gcli box claudex toon cybernest
)

echo "CPS v1.2.0 migration: moving tools from core to extensions"
echo ""

# Remove binaries from ~/shell/executables/
removed=0
for tool in "${MOVED_TOOLS[@]}"; do
  if [ -f "$EXEC_DIR/$tool" ]; then
    rm -f "$EXEC_DIR/$tool"
    echo "  removed binary: $tool"
    ((removed++))
  fi
done

if [ "$removed" -eq 0 ]; then
  echo "  no binaries to remove"
fi

# Clean entries from state.json
if [ -f "$STATE_FILE" ]; then
  if command -v jq &>/dev/null; then
    cleaned=0
    for tool in "${MOVED_TOOLS[@]}"; do
      if jq -e ".tools.\"$tool\"" "$STATE_FILE" &>/dev/null; then
        ((cleaned++))
      fi
    done

    if [ "$cleaned" -gt 0 ]; then
      # Build jq delete expression for all moved tools
      jq_expr="del("
      first=true
      for tool in "${MOVED_TOOLS[@]}"; do
        if [ "$first" = true ]; then
          jq_expr+=".tools.\"$tool\""
          first=false
        else
          jq_expr+=", .tools.\"$tool\""
        fi
      done
      jq_expr+=")"

      jq "$jq_expr" "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
      echo ""
      echo "  cleaned $cleaned entries from state.json"
    else
      echo ""
      echo "  no state entries to clean"
    fi
  else
    echo ""
    echo "  warning: jq not found, skipping state.json cleanup"
    echo "  install jq and re-run, or manually remove these tool entries from $STATE_FILE"
  fi
else
  echo ""
  echo "  no state.json found, skipping"
fi

echo ""
echo "done! reinstall with:"
echo "  cps extend security"
echo "  cps extend cloud"
echo "  cps extend appsec"
echo "  cps extend private  (requires --gh-token for toon, cybernest)"

# Duplicate PATH entries collapse to the first one, so sourcing this again from cps.zsh
# re-fronts the CPS paths after macOS path_helper has reordered them
typeset -U path PATH fpath FPATH

# --- Homebrew ---
export HOMEBREW_NO_AUTO_UPDATE=1
[ -f "$HOME/shell/env/brew.zsh" ] && source "$HOME/shell/env/brew.zsh"

# --- PATH ---
# custom-bin first so user-dropped binaries win over CPS-managed ones on name collision
export PATH="$HOME/shell/custom-bin:$HOME/shell/extensions:$HOME/.local/bin:$PATH"

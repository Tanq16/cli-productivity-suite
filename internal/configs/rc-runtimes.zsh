# --- Language Runtimes ---
export GOROOT="$HOME/shell/go-sdk"
export GOPATH="$HOME/shell/go"
export RUSTUP_HOME="$HOME/shell/rust/.rustup"
export CARGO_HOME="$HOME/shell/rust/.cargo"
export FNM_DIR="$HOME/shell/fnm"
export VIRTUAL_ENV="$HOME/shell/py-default"
export UV_TOOL_DIR="$HOME/shell/uv-tools"
export UV_TOOL_BIN_DIR="$HOME/shell/uv-tool-executables"

# --- PATH additions ---
export PATH="$UV_TOOL_BIN_DIR:$GOROOT/bin:$GOPATH/bin:$CARGO_HOME/bin:$VIRTUAL_ENV/bin:$PATH"

# --- Completions ---
[ -f "$HOME/shell/completions/uv.zsh" ] && source "$HOME/shell/completions/uv.zsh"
[ -f "$HOME/shell/completions/fnm.zsh" ] && source "$HOME/shell/completions/fnm.zsh"
command -v fnm &>/dev/null && eval "$(fnm env)"

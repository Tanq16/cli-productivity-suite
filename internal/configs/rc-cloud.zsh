# --- gcloud ---
if [ -n "$HOMEBREW_PREFIX" ]; then
  [ -f "$HOMEBREW_PREFIX/share/google-cloud-sdk/path.zsh.inc" ] && source "$HOMEBREW_PREFIX/share/google-cloud-sdk/path.zsh.inc"
  [ -f "$HOMEBREW_PREFIX/share/google-cloud-sdk/completion.zsh.inc" ] && source "$HOMEBREW_PREFIX/share/google-cloud-sdk/completion.zsh.inc"
fi

# --- AWS CLI ---
alias awsn='aws --no-cli-pager'
aws() { unset -f aws; autoload -Uz bashcompinit && bashcompinit; complete -C 'aws_completer' aws; command aws "$@"; }

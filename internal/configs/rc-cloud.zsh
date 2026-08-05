# --- gcloud ---
if [ -n "$HOMEBREW_PREFIX" ]; then
  [ -f "$HOMEBREW_PREFIX/share/google-cloud-sdk/path.zsh.inc" ] && source "$HOMEBREW_PREFIX/share/google-cloud-sdk/path.zsh.inc"
  # completion.zsh.inc costs ~30ms of a ~75ms startup; defer it to first use like aws below
  if [ -f "$HOMEBREW_PREFIX/share/google-cloud-sdk/completion.zsh.inc" ]; then
    gcloud() { unset -f gcloud; source "$HOMEBREW_PREFIX/share/google-cloud-sdk/completion.zsh.inc"; command gcloud "$@"; }
  fi
fi

# --- AWS CLI ---
alias awsn='aws --no-cli-pager'
aws() { unset -f aws; autoload -Uz bashcompinit && bashcompinit; complete -C 'aws_completer' aws; command aws "$@"; }

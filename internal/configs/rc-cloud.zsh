# --- gcloud ---
[ -f "$HOME/shell/gcloud-sdk/path.zsh.inc" ] && source "$HOME/shell/gcloud-sdk/path.zsh.inc"
[ -f "$HOME/shell/gcloud-sdk/completion.zsh.inc" ] && source "$HOME/shell/gcloud-sdk/completion.zsh.inc"

# --- AWS CLI ---
alias awsn='aws --no-cli-pager'
aws() { unset -f aws; autoload -Uz bashcompinit && bashcompinit; complete -C 'aws_completer' aws; command aws "$@"; }

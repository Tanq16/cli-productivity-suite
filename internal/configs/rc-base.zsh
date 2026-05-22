# --- Homebrew ---
export HOMEBREW_NO_AUTO_UPDATE=1
[ -f "$HOME/shell/env/brew.zsh" ] && source "$HOME/shell/env/brew.zsh"

# --- Zsh plugins (sourced directly, no framework) ---
ZSH_PLUGINS="$HOME/shell/plugins"
[ -f "$ZSH_PLUGINS/zsh-autosuggestions/zsh-autosuggestions.zsh" ] && source "$ZSH_PLUGINS/zsh-autosuggestions/zsh-autosuggestions.zsh"

# --- History ---
HISTFILE="$HOME/.zsh_history"
HISTSIZE=10000
SAVEHIST=10000
setopt HIST_IGNORE_DUPS HIST_IGNORE_SPACE SHARE_HISTORY INC_APPEND_HISTORY

# --- Completion ---
[ -d "$HOME/.cache/zsh" ] || mkdir -p "$HOME/.cache/zsh"
autoload -Uz compinit
compinit -d "$HOME/.cache/zsh/zcompdump"

# syntax-highlighting must be sourced LAST per upstream README
[ -f "$ZSH_PLUGINS/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh" ] && source "$ZSH_PLUGINS/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"

# --- PATH ---
# custom-bin first so user-dropped binaries win over CPS-managed ones on name collision
export PATH="$HOME/shell/custom-bin:$HOME/shell/extensions:$HOME/shell/executables:$HOME/.local/bin:$PATH"

# --- Starship prompt ---
command -v starship &>/dev/null && eval "$(starship init zsh)"

# --- FZF ---
[ -f "$HOME/shell/completions/fzf.zsh" ] && source "$HOME/shell/completions/fzf.zsh"
export FZF_DEFAULT_OPTS="
--layout=reverse
--info=inline
--height=95%
--multi
--preview '([[ -f {} ]] && (bat --style=numbers --color=always {} || cat {})) || ([[ -d {} ]] && (tree -C {} | less)) || echo {} 2>/dev/null | head -200'
--bind=ctrl-k:preview-down
--bind=ctrl-j:preview-up
"

export FZF_DEFAULT_COMMAND='fd --type f --hidden --exclude .git --exclude node_modules --exclude Library'
export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"

# --- Environment ---
export BAT_PAGER=''
export HISTCONTROL=ignoreboth
export LESS_TERMCAP_mb=$'\e[1;32m'
export LESS_TERMCAP_md=$'\e[1;32m'
export LESS_TERMCAP_me=$'\e[0m'
export LESS_TERMCAP_se=$'\e[0m'
export LESS_TERMCAP_so=$'\e[01;33m'
export LESS_TERMCAP_ue=$'\e[0m'
export LESS_TERMCAP_us=$'\e[1;4;31m'

# --- Aliases ---
alias f='fzf'
alias gitn='git --no-pager'
alias tt='tmux -u new -s default'
alias t='tmux -u a -t default'
alias tree='lsd --tree'
alias a=anbu
alias ts='tmux -u new -s'
alias ta='tmux -u a -t'
alias tls='tmux list-sessions'
alias vim=nvim
alias c=clear
alias l='lsd -l'
alias la='lsd -la'
alias sshide='ssh -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null"'
alias sessionrec='script -f $HOME/session-$(date +"%d-%b-%y_%H-%M-%S").log'
alias dockernonerm='for i in $(docker images -f dangling=true -q); do docker image rm $i; done'

# --- Zoxide ---
[ -f "$HOME/shell/completions/zoxide.zsh" ] && source "$HOME/shell/completions/zoxide.zsh"

# --- Cursor ---
keymap-for-cursor-shape() { echo -ne '\e[5 q' }
zle-line-init() { keymap-for-cursor-shape }

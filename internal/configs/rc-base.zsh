# --- Homebrew ---
export HOMEBREW_NO_AUTO_UPDATE=1
[ -f "$HOME/shell/env/brew.zsh" ] && source "$HOME/shell/env/brew.zsh"

# --- Oh My Zsh ---
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME="spaceship"
plugins=(zsh-autosuggestions zsh-syntax-highlighting)
source $ZSH/oh-my-zsh.sh

# --- PATH ---
export PATH="$HOME/shell/extensions:$HOME/shell/executables:$HOME/.local/bin:$PATH"

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
export SPACESHIP_TIME_SHOW=true
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

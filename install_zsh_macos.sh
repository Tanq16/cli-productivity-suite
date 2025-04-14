#!/bin/sh

printf "Installing several packages; may take a few minutes. Hang tight.\n\n"

# brew installs
brew install tree tmux jq bat fd neovim neofetch ripgrep lsd wget git 1>/dev/null 2>/dev/null
brew install --cask nikitabobko/tap/aerospace 1>/dev/null 2>/dev/null
printf '[+] Installed brew packages'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i '' -e "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -i '' -e "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
printf '\r\033[K'
printf '[+] Installed OMZ and plugins'

# Fix slow ZSH paste
sed -i '' -e "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i '' -e "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i '' -e "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i '' -e "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh

# NeoVIM setup
rm -rf ~/.vim* 1>/dev/null 2>/dev/null
rm -rf ~/.config/nvim 2>/dev/null
rm -rf ~/.local/share/nvim 2>/dev/null
git clone https://github.com/NvChad/starter ~/.config/nvim 1>/dev/null 2>/dev/null
sed -i '' -e 's/theme =.*/theme = "catppuccin", transparency = true, --/' ~/.config/nvim/lua/chadrc.lua
nvim --headless -c 'quitall' 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installed NeoVIM'

# TMUX & FZF
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 1>/dev/null 2>/dev/null
sed -i '' -e 's/bind-key -T copy-mode MouseDragEnd1Pane send -X copy-selection/bind-key -T copy-mode MouseDragEnd1Pane send-keys -X copy-pipe "pbcopy"/' ./tmuxconf # fix copy on select in tmux for macOS
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 1>/dev/null 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installed TMUX & FZF'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/macos.rcfile 2>/dev/null
cat ~/.zshrc >> ./temptemp
cat ./macos.rcfile >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./macos.rcfile 1>/dev/null 2>/dev/null
printf '\r\033[L'
printf '[+] Setup RC files'

printf "\n\nNew shell spawns in 10 seconds; quit the terminal for everything to take effect then start again.\n"
sleep 10
exec zsh -l

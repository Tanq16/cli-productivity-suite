#!/bin/sh

printf "\nIf you have some other vim config installed, press ⌃+c now and remove that.\nAfter setting up appropriately, start the script again. Sleeping for 20 seconds."
for i in $(seq 20); do printf '.'; sleep 1; done
printf "\n\nInstalling several packages; may take a few minutes. Hang tight.\n"

# brew installs
brew install tree tmux jq bat fd neovim neofetch ripgrep lsd wget git 1>/dev/null 2>/dev/null
brew install --cask nikitabobko/tap/aerospace 1>/dev/null 2>/dev/null
printf '.'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i '' -e "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -i '' -e "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
printf '.'

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
printf '.'
nvim --headless -c 'quitall' 1>/dev/null 2>/dev/null
printf '.'

# TMUX & FZF
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 1>/dev/null 2>/dev/null
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null
printf '.'
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 1>/dev/null 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null
printf '.'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/macos.rcfile 2>/dev/null
cat ~/.zshrc >> ./temptemp
cat ./macos.rcfile >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./macos.rcfile 1>/dev/null 2>/dev/null

printf "\n\nIf you don't see shapes properly after this, make sure to install the font properly (Read the README)\n\nNOTE: After the new shell spawns, quit the terminal app or SSH session for everything to take effect then start again\n\nStarting in 10 seconds."
for i in $(seq 10); do printf '.'; sleep 1; done;
exec zsh -l

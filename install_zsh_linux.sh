#!/bin/sh

printf "\nIf you have some other vim config installed, press ⌃+c now and remove that.\nAfter setting up appropriately, start the script again. Sleeping for 20 seconds."
for i in $(seq 20); do printf '.'; sleep 1; done
printf "\n\nInstalling several packages; may take a few minutes. Hang tight.\n"

# apt installs
sudo apt update -y 1>/dev/null 2>/dev/null && sudo apt install -y tar wget tree tmux ripgrep jq ninja-build gettext neofetch make cmake unzip curl git file gcc bat fd-find 1>/dev/null 2>/dev/null
printf '.'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
printf '.'

# Fix slow ZSH paste
sed -i "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh

# NeoVIM setup
rm -rf ~/.vim* 1>/dev/null 2>/dev/null
rm -rf ~/.config/nvim 2>/dev/null
rm -rf ~/.local/share/nvim 2>/dev/null
sudo apt remove vim neovim -y 1>/dev/null 2>/dev/null
wget https://github.com/neovim/neovim/archive/refs/tags/stable.tar.gz 1>/dev/null 2>/dev/null
tar -xvf stable.tar.gz 1>/dev/null 2>/dev/null
cd neovim-stable
printf '.'
make CMAKE_BUILD_TYPE=RelWithDebInfo 1>/dev/null 2>/dev/null
printf '.'
cd build
cpack -G DEB 1>/dev/null 2>/dev/null
printf '.'
sudo apt install ./nvim-linux64.deb 1>/dev/null 2>/dev/null
cd ../.. && rm -rf stable.tar.gz neovim-stable 1>/dev/null 2>/dev/null
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

# colored ls - lsd
a=$(curl -s https://api.github.com/repos/lsd-rs/lsd/releases/latest | grep -E "browser_download_url.*" | grep -i "linux" | grep -i "x86_64" | grep -i "gnu" | cut -d '"' -f4)
wget "$a" -O test.tar.gz 1>/dev/null 2>/dev/null
tar -xzf test.tar.gz 1>/dev/null 2>/dev/null
sudo mv lsd-*/lsd /usr/bin/lsd
rm -rf lsd-* test.tar.gz
printf '.'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/linux.rcfile 1>/dev/null 2>/dev/null
cat ~/.zshrc >> ./temptemp
cat ./linux.rcfile >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./linux.rcfile 1>/dev/null 2>/dev/null

printf "\n\nIf you don't see shapes properly after this, make sure to install the font properly (Read the README)\n\nNOTE: After the new shell spawns, quit the terminal app or SSH session for everything to take effect then start again\n\nStarting in 10 seconds."
for i in $(seq 10); do printf '.'; sleep 1; done;
exec zsh -l

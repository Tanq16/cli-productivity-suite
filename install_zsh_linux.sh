#!/bin/sh

printf "Installing several packages; may take a few minutes. Hang tight.\n\n"

# apt installs
sudo apt update -y 1>/dev/null 2>/dev/null && sudo apt install -y tar wget tree tmux ripgrep jq ninja-build gettext neofetch make cmake unzip curl git file gcc bat fd-find 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installed apt packages'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
printf '\r\033[K'
printf '[+] Installed OMZ and plugins'

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
printf '\r\033[K'
printf '[+] Building NeoVIM'
make CMAKE_BUILD_TYPE=RelWithDebInfo 1>/dev/null 2>/dev/null
cd build
cpack -G DEB 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installing NeoVIM'
sudo apt install ./nvim-linux-x86_64.deb 1>/dev/null 2>/dev/null
sudo apt install ./nvim-linux-arm64.deb 1>/dev/null 2>/dev/null
cd ../.. && rm -rf stable.tar.gz neovim-stable 1>/dev/null 2>/dev/null
git clone https://github.com/NvChad/starter ~/.config/nvim 1>/dev/null 2>/dev/null
sed -i 's/theme =.*/theme = "catppuccin", transparency = true,/' ~/.config/nvim/lua/chadrc.lua
nvim --headless -c 'quitall' 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installed NeoVIM'

# TMUX & FZF
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 1>/dev/null 2>/dev/null
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 1>/dev/null 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null
printf '\r\033[K'
printf '[+] Installed TMUX & FZF'

# colored ls - lsd
a=$(curl -s https://api.github.com/repos/lsd-rs/lsd/releases/latest | grep -E "browser_download_url.*" | grep -i "linux" | grep -i "x86_64" | grep -i "gnu" | cut -d '"' -f4)
wget "$a" -O test.tar.gz 1>/dev/null 2>/dev/null
tar -xzf test.tar.gz 1>/dev/null 2>/dev/null
sudo mv lsd-*/lsd /usr/bin/lsd
rm -rf lsd-* test.tar.gz
printf '\r\033[K'
printf '[+] Installed colored ls - lsd'

# install go
a=$(curl -s https://go.dev/dl/ | grep -oE "(/dl/go[\.0-9]{2,7}\.linux-arm64\.tar\.gz)" | head -n 1) 1>/dev/null 2>/dev/null && \
b=$(curl -s https://go.dev/dl/ | grep -oE "(/dl/go[\.0-9]{2,7}\.linux-amd64\.tar\.gz)" | head -n 1) 1>/dev/null 2>/dev/null && \
if [ "$(uname -m)" = "aarch64" ]; then wget "https://golang.org$a" 1>/dev/null 2>/dev/null; else wget "https://golang.org$b" 1>/dev/null 2>/dev/null; fi && \
if [ "$(uname -m)" = "aarch64" ]; then c=$(echo $a | cut -d "/" -f3); else c=$(echo $b | cut -d "/" -f3); fi && \
sudo tar -C /usr/local -xzf "$c" 1>/dev/null 2>/dev/null && \
rm "$c"
printf '\r\033[K'
printf '[+] Installed Go'

# personal executables
mkdir -p $HOME/shell/executables
/usr/local/go/bin/go install github.com/tanq16/ai-context@latest 1>/dev/null 2>/dev/null
/usr/local/go/bin/go install github.com/tanq16/nottif@latest 1>/dev/null 2>/dev/null
/usr/local/go/bin/go install github.com/tanq16/danzo@latest 1>/dev/null 2>/dev/null
mv $HOME/go/bin/ai-context $HOME/shell/executables/ai-context
mv $HOME/go/bin/nottif $HOME/shell/executables/nottif
mv $HOME/go/bin/danzo $HOME/shell/executables/danzo
printf '\r\033[K'
printf '[+] Installed personal executables'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/linux.rcfile 1>/dev/null 2>/dev/null
cat ~/.zshrc >> ./temptemp
cat ./linux.rcfile >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./linux.rcfile 1>/dev/null 2>/dev/null

printf "\n\nNew shell spawns in 10 seconds; quit the terminal app for everything to take effect then start again.\n"
sleep 10
exec zsh -l

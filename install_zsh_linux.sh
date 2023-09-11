#!/bin/sh

echo ""
echo "If you have some other vim config installed, press ⌃+c now and remove that."
echo "After setting up appropriately, start the script again. Sleeping for 20 seconds!"
for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do echo -n '.'; sleep 1; done; echo ""

echo "\n"
echo "Installing several packages; may take a few minutes. Hang tight!\n"

# apt installs
sudo apt update -y 1>/dev/null 2>/dev/null && sudo apt install -y tar wget tree tmux jq ninja-build gettext make cmake unzip curl git file gcc bat fd-find 1>/dev/null 2>/dev/null
echo -n '.'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
echo -n '.'

# NeoVIM setup
rm -rf ~/.vim* 1>/dev/null 2>/dev/null
rm -rf ~/.config/nvim 2>/dev/null
rm -rf ~/.local/share/nvim 2>/dev/null
sudo apt remove vim neovim -y 1>/dev/null 2>/dev/null
wget https://github.com/neovim/neovim/archive/refs/tags/stable.tar.gz 1>/dev/null 2>/dev/null
tar -xvf stable.tar.gz 1>/dev/null 2>/dev/null
cd neovim-stable
echo -n '.'
make CMAKE_BUILD_TYPE=RelWithDebInfo 1>/dev/null 2>/dev/null
echo -n '.'
cd build
cpack -G DEB 1>/dev/null 2>/dev/null
echo -n '.'
sudo apt install ./nvim-linux64.deb 1>/dev/null 2>/dev/null
cd ../.. && rm -rf stable.tar.gz neovim-stable 1>/dev/null 2>/dev/null
git clone https://github.com/NvChad/NvChad ~/.config/nvim --depth 1 1>/dev/null 2>/dev/null
sed -i "s/local input =/local input = \"N\" --/" ~/.config/nvim/lua/core/bootstrap.lua
sed -i "s/dofile(vim.g/vim.cmd([[ set guicursor= ]])\ndofile(vim.g/" ~/.config/nvim/init.lua
echo -n '.'
nvim --headless -c 'quitall' 1>/dev/null 2>/dev/null
echo -n '.'

# TMUX & FZF
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 1>/dev/null 2>/dev/null
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null
echo -n '.'
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 1>/dev/null 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null
echo -n '.'

# colored ls - lsd
if [ $(uname -p) != "x86_64" ]
then
    a=$(curl -L -s https://github.com/Peltoche/lsd/releases/latest | grep -oE "tag.+\"" | cut -d '/' -f2 | grep -vE "^[^0-9]" | cut -d "\"" -f1 | head -n 1) && \
    wget "https://github.com/Peltoche/lsd/releases/download/$a/lsd_""$a""_arm64.deb" 1>/dev/null 2>/dev/null && \
    sudo apt install -y "./lsd_""$a""_arm64.deb" 1>/dev/null 2>/dev/null && \
    rm "lsd_""$a""_arm64.deb"
else
    a=$(curl -L -s https://github.com/Peltoche/lsd/releases/latest | grep -oE "tag.+\"" | cut -d '/' -f2 | grep -vE "^[^0-9]" | cut -d "\"" -f1 | head -n 1) && \
    wget "https://github.com/Peltoche/lsd/releases/download/$a/lsd_""$a""_amd64.deb" 1>/dev/null 2>/dev/null && \
    sudo apt install -y "./lsd_""$a""_amd64.deb" 1>/dev/null 2>/dev/null && \
    rm "lsd_""$a""_amd64.deb"
fi
echo -n '.'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 1>/dev/null 2>/dev/null
cat ~/.zshrc >> ./temptemp
cat ./add_to_rc >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./add_to_rc 1>/dev/null 2>/dev/null

echo "\n\n"
echo "If you don't see shapes properly after this, make sure to install the font properly (Read the README)"
echo ""
echo "NOTE: After the new shell spawns, quit the terminal app or SSH session for everything to take effect then start again"
echo ""
echo "Starting in 10 seconds!"
for i in 1 2 3 4 5 6 7 8 9 10; do echo -n '.'; sleep 1; done; echo ""
exec zsh -l

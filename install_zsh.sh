#!/bin/sh

echo "If you have some other vim config installed, press ⌃+c now and remove that."
echo "After setting up appropriately, start the script again. Sleeping for 20 seconds!"
for i in {1..20}
do
    echo -n '.'
    sleep 1
done

echo ""
echo ""

echo "Insitializing..... may take a few seeconds."
if [ $(uname -s) != "Darwin" ]
then
    sudo apt install -y tree tmux jq 1>/dev/null 2>/dev/null
else
    brew install tree tmux jq 1>/dev/null 2>/dev/null
fi

git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
if [ $(uname -s) != "Darwin" ]
then
    sed -i "s/robbyrussell/spaceship/" ~/.zshrc
else
    sed -ie "s/robbyrussell/spaceship/" ~/.zshrc
fi

echo "Custom shell installed. Adding plugins!"
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null

if [ $(uname -s) != "Darwin" ]
then
    sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
else
    sed -ie "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
fi

echo "Installing NeoVIM, Bat and fd-Find"
rm -rf ~/.vim* 1>/dev/null 2>/dev/null
if [ $(uname -s) != "Darwin" ]
then
    sudo apt install bat fd-find -y 1>/dev/null 2>/dev/null
    if [ $(uname -p) != "x86_64" ]
    then
        echo "Assuming NeoVIM installed for ARM linux..."
    else
        wget https://github.com/Tanq16/cli-productivity-suite/releases/download/x86_64-deb/nvim-linux64.deb 1>/dev/null 2>/dev/null
        sudo apt install ./nvim-linux64.deb 1>/dev/null 2>/dev/null
        rm nvim-linux64.deb
    fi
else
    brew install bat fd neovim 1>/dev/null 2>/dev/null
fi

echo "Installing Tmux"
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 2>/dev/null
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null

echo "Installing fuzzy finder"
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null

echo "Installing colored ls"
if [ $(uname -s) != "Darwin" ]
then
    if [ $(uname -p) != "x86_64" ]
    then
        a=$(curl -L -s https://github.com/Peltoche/lsd/releases/latest | grep -oE "tag.+\"" | cut -d '/' -f2 | grep -vE "^[^0-9]" | cut -d "\"" -f1 | head -n 1) && \
        wget "https://github.com/Peltoche/lsd/releases/download/$a/lsd_""$a""_arm64.deb" 2>/dev/null 1>/dev/null && \
        sudo apt install -y "./lsd_""$a""_arm64.deb" 2>/dev/null 1>/dev/null && \
        rm "lsd_""$a""_arm64.deb"
    else
        a=$(curl -L -s https://github.com/Peltoche/lsd/releases/latest | grep -oE "tag.+\"" | cut -d '/' -f2 | grep -vE "^[^0-9]" | cut -d "\"" -f1 | head -n 1) && \
        wget "https://github.com/Peltoche/lsd/releases/download/$a/lsd_""$a""_amd64.deb" 2>/dev/null 1>/dev/null && \
        sudo apt install -y "./lsd_""$a""_amd64.deb" 2>/dev/null 1>/dev/null && \
        rm "lsd_""$a""_amd64.deb"
    fi
else
    brew install lsd 1>/dev/null 2>/dev/null
fi

echo "Installing NvChad for NeoVIM."
rm -rf ~/.config/nvim 2>/dev/null
rm -rf ~/.local/share/nvim 2>/dev/null
git clone https://github.com/NvChad/NvChad ~/.config/nvim --depth 1 1>/dev/null 2>/dev/null
if [ $(uname -s) != "Darwin" ]
then
    sed -i "s/local input =/local input = N --/" ~/.config/nvim/lua/core/bootstrap.lua
    sed -i "s/dofile(vim.g/vim.cmd([[ set guicursor= ]])\ndofile(vim.g/" ~/.config/nvim/init.lua
else
    sed -ie "s/local input =/local input = N --/" ~/.config/nvim/lua/core/bootstrap.lua
    sed -ie "s/dofile(vim.g/vim.cmd([[ set guicursor= ]])\ndofile(vim.g/" ~/.config/nvim/init.lua
fi
nvim --headless -c 'quitall'

wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
if [ $(uname -s) != "Darwin" ]
then
    echo ""
else
    sed -ie "s/fdfind/fd/" ./add_to_rc
    sed -ie "s/batcat/bat/" ./add_to_rc
    sed -ie "s/alias bat=/\# alias bat=/" ./add_to_rc
    sed -ie "s/alias ip4/alias ip4='ifconfig | grep -oE \"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\" | grep -v \"127.0.0.1\" | grep -vE \".+\.255$\"' #/" ./add_to_rc
fi
cat ./add_to_rc >> ./temptemp
cat temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./add_to_rc ./add_to_rce ./.zshrce 2>/dev/null

echo "\n\n\n\n\n\n\n\n"
echo "If you don't see shapes properly, make sure to install the font properly (Read the README)"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting in 5 seconds."
echo "Ater the configuration, restart the shell (terminal app for ease) for everything to take effect"
sleep 5
exec zsh -l

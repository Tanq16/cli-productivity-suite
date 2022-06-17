#!/bin/sh

echo "Insitializing..... may take a few seeconds."
if [ $(uname -s) != "Darwin" ]
then
    sudo apt install -y tree tmux 1>/dev/null 2>/dev/null
else
    brew install tree tmux 1>/dev/null 2>/dev/null
fi

# echo "Setting ZSH to default shell :: Please enter your password."
# chsh -s /usr/bin/zsh $USER

# git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
git clone https://github.com/spaceship-prompt/spaceship-prompt.git "~/.oh-my-zsh/custom/themes/spaceship-prompt" --depth=1
ln -s "~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme" "~/.oh-my-zsh/custom/themes/spaceship.zsh-theme"
if [ $(uname -s) != "Darwin" ]
then
    # sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc
    sed -i "s/robbyrussell/spaceship/" ~/.zshrc
else
    sed -ie "s/robbyrussell/spaceship/" ~/.zshrc
fi

echo "Custom shell installed."
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null

if [ $(uname -s) != "Darwin" ]
then
    sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
else
    sed -ie "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
fi

echo "Installing Bat and fd-Find"
if [ $(uname -s) != "Darwin" ]
then
    sudo apt install bat fd-find -y 1>/dev/null 2>/dev/null
else
    brew install bat fd 1>/dev/null 2>/dev/null
fi

echo "Installing Tmux - On first start, press Prefix (ctrl+b) then press shift-i to install plugins properly."
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 2>/dev/null
mv tmuxconf ~/.tmux.conf

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

curl -sLf https://spacevim.org/install.sh | bash
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/spacevim_config 2>/dev/null
mkdir ~/.SpaceVim.d
mv spacevim_config ~/.SpaceVim.d/

wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
if [ $(uname -s) != "Darwin" ]
then
    echo "finishing up"
else
    sed -ie "s/fdfind/fd/" ./add_to_rc
    sed -ie "s/batcat/bat/" ./add_to_rc
    sed -ie "s/alias bat=/\# alias bat=/" ./add_to_rc
fi
cat add_to_rc >> ~/.zshrc
rm add_to_rc
cp .zshrc temptemp
cat temptemp | grep -vE "^#" | grep -vE "^$" > .zshrc
rm temptemp

echo "\n\n\n\n\n\n\n\n"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "If you don't see shapes properly, install powerline fonts (Read the README)"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting powerlevel config in 7 seconds."
echo "Ater the configuration, close all shell instances and restart shell for all plugins to take effect"
sleep 7
exec zsh -l

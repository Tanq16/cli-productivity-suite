#!/bin/sh

echo "Insitializing..... may take a few seeconds."
if [ $(uname -s) = "Darwin" ]
then
    sudo apt install -y tree sshpass tmux 1>/dev/null 2>/dev/null
else
    brew install tree tmux 1>/dev/null 2>/dev/null
fi

# echo "Setting ZSH to default shell :: Please enter your password."
# chsh -s /usr/bin/zsh $USER

git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc

echo "Custom shell installed."
echo "Installing Auto-suggestions"
git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
echo "Installing Syntax highlighting"
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null

sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc

echo "Installing Bat and fd-Find"
if [ $(uname -s) = "Darwin" ]
then
    sudo apt install bat -y 1>/dev/null 2>/dev/null
    sudo apt install fd-find -y 1>/dev/null 2>/dev/null
else
    brew install vim bat fd 1>/dev/null 2>/dev/null
fi

echo "Installing Tmux - Upon first start, press Prefix then type :source-file ~/.tmux.conf, then press Prefix->I to install plugins properly."
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 2>/dev/null
git clone https://github.com/jimeh/tmux-themepack.git ~/.tmux-themepack
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 2>/dev/null
mv tmuxconf ~/.tmux.conf

echo "Installing fuzzy finder"
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null

echo "Installing colored ls"
if [ $(uname -s) = "Darwin" ]
then
    wget https://github.com/Peltoche/lsd/releases/download/0.20.1/lsd_0.20.1_amd64.deb 2>/dev/null
    apt install -y ./lsd_0.20.1_amd64.deb 1>/dev/null 2>/dev/null && rm lsd_0.20.1_amd64.deb
else
    brew install lsd 1>/dev/null 2>/dev/null
fi

wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
cat add_to_rc >> ~/.zshrc
rm add_to_rc
cp .zshrc temptemp
cat temptemp | grep -vE "^#" | grep -vE "^$" > .zshrc
rm temptemp

echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "If you don't see shapes properly, install powerline fonts. The recommended font is given in the README of the repo for this script."
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting powerlevel config in 10 seconds. Ater the configuration, please close all shell instances and restart the shell for fuzzy search and zsh plugins to take effect......."
sleep 10
exec zsh -l

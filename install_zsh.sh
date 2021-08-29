#!/bin/sh

echo "Insitializing..... may take a few seeconds."
sudo apt install -y tree sshpass tmux 1>/dev/null 2>/dev/null

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

echo "Installing Awesome color scheme"
git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
git clone https://github.com/arcticicestudio/nord-dircolors.git 2>/dev/null
cp nord-dircolors/src/dir_colors ~/.oh-my-zsh/nord.dircolors
rm -rf nord-dircolors/
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/

echo "Installing Bat and fd-Find"
sudo apt install bat -y 1>/dev/null 2>/dev/null
sudo apt install fd-find -y 1>/dev/null 2>/dev/null

echo "Installing Tmux - Upon first start, press Prefix then type :source-file ~/.tmux.conf, then press Prefix->I to install plugins properly."
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 2>/dev/null
git clone https://github.com/jimeh/tmux-themepack.git ~/.tmux-themepack
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 2>/dev/null
mv tmuxconf ~/.tmux.conf

echo "Installing fuzzy finder"
sudo apt install fzf
echo "if [[ ! "$PATH" == */home/tanq/.fzf/bin* ]]; then" > ~/.fzf.zsh
echo '  export PATH="${PATH:+${PATH}:}/home/tanq/.fzf/bin"' >> ~/.fzf.zsh
echo 'fi' >> ~/.fzf.zsh
echo '[[ $- == *i* ]] && source "/home/tanq/.fzf/shell/completion.zsh" 2> /dev/null' >> ~/.fzf.zsh
echo 'source "/home/tanq/.fzf/shell/key-bindings.zsh"' >> ~/.fzf.zsh

echo "Installing colored ls"
wget https://github.com/Peltoche/lsd/releases/download/0.20.1/lsd_0.20.1_amd64.deb 2>/dev/null
apt install -y ./lsd_0.20.1_amd64.deb 1>/dev/null 2>/dev/null && rm lsd_0.20.1_amd64.deb

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

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
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null

wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
cat add_to_rc >> ~/.zshrc
rm add_to_rc
cp .zshrc temptemp
cat temptemp | grep -vE "^#" | grep -vE "^$" > .zshrc

echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "If you don't see shapes properly, install powerline fonts. The recommended font is given in the README of the repo for this script."
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting powerlevel config in 10 seconds. Ater the configuration, please close all shell instances and restart the shell for fuzzy search and zsh plugins to take effect......."
sleep 10
exec zsh -l

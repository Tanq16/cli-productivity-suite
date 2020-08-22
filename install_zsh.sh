#!/bin/sh

echo "Installing ZSH, wget and git. This may take 3-4 minutes depending on network/processor/storage."
sudo apt install -y zsh wget git 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh --unattended
rm install.sh

echo "Setting ZSH to default shell :: Please enter your password."
chsh -s /usr/bin/zsh $USER

git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc

echo "Custom shell installed."
echo "Installing Auto-suggestions"
sudo apt install -y zsh-autosuggestions 1>/dev/null 2>/dev/null
echo "Installing Syntax highlighting"
sudo apt install -y zsh-syntax-highlighting 1>/dev/null 2>/dev/null

sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
echo "alias c=clear" >> ~/.zshrc
echo "alias l='ls -l'" >> ~/.zshrc
echo "alias la='ls -la'" >> ~/.zshrc
echo "export BAT_PAGER=''" >> ~/.zshrc

echo "Installing Awesome color scheme"
git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/
echo "eval \`dircolors ~/.oh-my-zsh/dircolors.256dark\`" >> ~/.zshrc

wget https://github.com/sharkdp/bat/releases/download/v0.11.0/bat_0.11.0_amd64.deb 2>/dev/null
sudo dpkg -i bat_0.11.0_amd64.deb 2>/dev/null
rm bat_0.11.0_amd64.deb

echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "If you don't see shapes properly, install powerline fonts. The recommended font is given in the README of the repo for this script."
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting powerlevel config in 10 seconds ..."
sleep 10
exec zsh -l

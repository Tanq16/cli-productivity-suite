#!/bin/sh

echo "Installing ZSH, wget and git. This may take 3-4 minutes depending on network/processor/storage."
sudo apt install -y zsh wget git 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh
sh install.sh --unattended

echo "Setting ZSH to default shell :: Please enter your password."
chsh -s /usr/bin/zsh $USER

git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc

echo "Custom shell installed."
echo "Installing Auto-suggestions"
git clone https://github.com/zsh-users/zsh-autosuggestions.git $ZSH_CUSTOM/plugins/zsh-autosuggestions 2>/dev/null
echo "Installing Syntax highlighting"
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git $ZSH_CUSTOM/plugins/zsh-syntax-highlighting 2>/dev/null

sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
echo "alias c=clear" >> ~/.zshrc
echo "alias l='ls -l'" >> ~/.zshrc
echo "alias la='ls -la'" >> ~/.zshrc

echo "Installing Awesome color scheme"
git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/
echo "eval \`dircolors ~/.oh-my-zsh/dircolors.256dark\`" >> ~/.zshrc

echo "If you chose Powerlevel10k installation, the installation options"
exec zsh -l

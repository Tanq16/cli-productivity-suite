#!/bin/sh

sudo apt install -y zsh curl git
curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh > install_hahanet.sh
sh install_hahanet.sh --unattended

echo """______________________________________________________________________________

Choose an installation:
1. OhMyZsh (Recommended, good for all)
2. Powerlevel10k (Fast alternative, with more features that need installation. Good for heavy customization)
______________________________________________________________________________
"""

read choice
if [ $choice -eq 2 ]
then
        git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k
        sed -i "s/robbyrussel/powerlevel10k\/powerlevel10k/" ~/.zshrc
else
        sed -i "s/robbyrussel/agnoster/" ~/.zshrc
fi

sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc

echo "alias c=clear" >> .zshrc
echo "alias l='ls -l'" >> .zshrc
echo "alias la='ls -la'" >> .zshrc

git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/
echo "eval \`dircolors ~/.oh-my-zsh/dircolors.256dark\`" >> .zshrc

source ~/.zshrc

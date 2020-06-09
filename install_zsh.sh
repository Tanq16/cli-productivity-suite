#!/bin/sh

sudo apt install zsh curl git
sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
git clone https://github.com/zsh-users/zsh-autosuggestions.git $ZSH_CUSTOM/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git $ZSH_CUSTOM/plugins/zsh-syntax-highlighting
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k
cat ~/.zshrc | sed -e "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" > temp_file_for_zsh_rc
cat temp_file_for_zsh_rc > ~/.zshrc
rm temp_file_for_zsh_rc
cat ~/.zshrc | sed -e "s/robbyrussel/powerlevel10k/powerlevel10k/" > temp_file_for_zsh_rc
cat temp_file_for_zsh_rc > ~/.zshrc
rm temp_file_for_zsh_rc

echo "alias c=clear" >> .zshrc
echo "alias l='ls -l'" >> .zshrc
echo "alias la='ls -la'" >> .zshrc

git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/
echo "eval \`dircolors ~/.oh-my-zsh/dircolors.256dark\`" >> .zshrc

source ~/.zshrc


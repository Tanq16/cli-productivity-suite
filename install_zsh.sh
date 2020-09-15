#!/bin/sh

echo "Installing ZSH, wget and git. This may take 3-4 minutes depending on network/processor/storage."
sudo apt install -y zsh wget git tree sshpass 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh --unattended
rm install.sh

echo "Setting ZSH to default shell :: Please enter your password."
chsh -s /usr/bin/zsh $USER

git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/themes/powerlevel10k 2>/dev/null
sed -i "s/robbyrussell/powerlevel10k\/powerlevel10k/" ~/.zshrc

echo "Custom shell installed."
echo "Installing Auto-suggestions"
git clone https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
echo "Installing Syntax highlighting"
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null

sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
echo "alias c=clear" >> ~/.zshrc
echo "alias l='ls -l'" >> ~/.zshrc
echo "alias la='ls -la'" >> ~/.zshrc
echo "export BAT_PAGER=''" >> ~/.zshrc
echo "alias bat=batcat" >> ~/.zshrc

echo "Installing Awesome color scheme"
git clone https://github.com/seebi/dircolors-solarized.git dirco_for_script_color_option 2>/dev/null
cp dirco_for_script_color_option/dircolors.256dark ~/.oh-my-zsh/dircolors.256dark
rm -rf dirco_for_script_color_option/
echo "eval \`dircolors ~/.oh-my-zsh/dircolors.256dark\`" >> ~/.zshrc

echo "Installing Bat and fast-update command"
mkdir -p ~/.custom_commands
touch ~/.custom_commands/update
echo "#!/usr/bin/zsh" >> ~/.custom_commands/update
echo "echo 'The install process requires 2-30 mins depending on size of updates and speed of drives and network.'" >> ~/.custom_commands/update
echo "echo 'Fetching repositories .....'" >> ~/.custom_commands/update
echo "a=$(sudo apt update 2>/dev/null)" >> ~/.custom_commands/update
echo "if [[ $a == *--upgradable* ]]" >> ~/.custom_commands/update
echo "then" >> ~/.custom_commands/update
echo "    echo 'Upgrading packages'" >> ~/.custom_commands/update
echo "    sudo apt upgrade -y 1>/dev/null 2>/dev/null" >> ~/.custom_commands/update
echo "    sudo apt dist-upgrade -y 1>/dev/null 2>/dev/null" >> ~/.custom_commands/update
echo "    echo 'Removing unused packages'" >> ~/.custom_commands/update
echo "    sudo apt autoremove -y 1>/dev/null 2>/dev/null" >> ~/.custom_commands/update
echo "    echo 'Clearing local repository'" >> ~/.custom_commands/update
echo "    sudo apt autoclean -y 1>/dev/null 2>/dev/null" >> ~/.custom_commands/update
echo "    echo 'System Updated !!'" >> ~/.custom_commands/update
echo "else" >> ~/.custom_commands/update
echo "    echo 'System upto date !!'" >> ~/.custom_commands/update
echo "fi" >> ~/.custom_commands/update
chmod +x ~/.custom_commands/update
echo "export PATH=$PATH:~/custom_commands/" >> ~/.zshrc

sudo apt install bat -y 1>/dev/null 2>/dev/null
sudo apt install fd-find -y 1>/dev/null 2>/dev/null

echo "Downloading fuzzy finder"
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null

echo 'export FZF_DEFAULT_OPTS="' >> ~/.zshrc
echo "--layout=reverse" >> ~/.zshrc
echo "--info=inline" >> ~/.zshrc
echo "--height=95%" >> ~/.zshrc
echo "--multi" >> ~/.zshrc
echo "--preview '([[ -f {}  ]] && (batcat --style=numbers --color=always {} || cat {})) || ([[ -d {}  ]] && (tree -C {} | less)) || echo {} 2> /dev/null | head -200'" >> ~/.zshrc
echo "--bind=ctrl-k:preview-down" >> ~/.zshrc
echo "--bind=ctrl-j:preview-up" >> ~/.zshrc
echo '"' >> ~/.zshrc
echo 'alias f=fzf' >> ~/.zshrc
echo "export FZF_DEFAULT_COMMAND='fdfind --follow --hidden'" >> ~/.zshrc
echo 'export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"' >> ~/.zshrc

echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "If you don't see shapes properly, install powerline fonts. The recommended font is given in the README of the repo for this script."
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo "Starting powerlevel config in 10 seconds. Ater the configuration, please close all shell instances and restart the shell for fuzzy search and zsh plugins to take effect......."
sleep 10
exec zsh -l

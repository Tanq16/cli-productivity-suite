#!/usr/bin/zsh
echo "Installing vim"
sudo apt install -y vim 1>/dev/null 2>/dev/null

echo "Installing Lightline"
git clone https://github.com/itchyny/lightline.vim ~/.vim/pack/plugins/start/lightline 2>/dev/null
echo "Installing Bracket Auto-Pair"
git clone https://github.com/jiangmiao/auto-pairs.git ~/auto_pairs_vim 2>/dev/null
mkdir -p ~/.vim/plugin
cp ~/auto_pairs_vim/plugin/auto-pairs.vim ~/.vim/plugin/
rm -rf ~/auto_pairs_vim

echo "Installing SuperTab"
wget https://www.vim.org/scripts/download_script.php\?src_id\=21752 2>/dev/null
mv 'download_script.php?src_id=21752' ~/.vim/supertab.vmb
vim -c 'so %' -c 'q' ~/.vim/supertab.vmb

echo "Installing nord colorscheme"
wget https://raw.githubusercontent.com/dylnmc/novum.vim/master/colors/novum.vim 2>/dev/null
git clone https://github.com/arcticicestudio/nord-vim.git 2>/dev/null
mkdir -p ~/.vim/colors
mv nord-vim/colors/nord.vim ~/.vim/colors/nord.vim
mv nord-vim/autoload/lightline/colorscheme/nord.vim ~/.vim/pack/plugins/start/lightline/autoload/lightline/colorscheme/
rm -rf nord-vim/

echo "Installing nerd commenter"
curl -fLo ~/.vim/plugin/NERD_Commenter.vim --create-dirs https://raw.githubusercontent.com/preservim/nerdcommenter/master/plugin/NERD_commenter.vim 2>/dev/null
curl -fLo ~/.vim/doc/NERD_Commenter.txt --create-dirs https://raw.githubusercontent.com/preservim/nerdcommenter/master/doc/NERD_commenter.txt 2>/dev/null

mv ~/.vimrc ~/.vimrc.old
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/.vimrcfile 2>/dev/null
mv .vimrcfile ~/.vimrc

sleep 2
echo "Done!"

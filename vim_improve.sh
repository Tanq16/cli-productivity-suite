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
mkdir -p ~/.vim/pack/plugins/start
git clone --depth=1 https://github.com/ervandew/supertab.git ~/.vim/pack/plugins/start/supertab

echo "Installing colorscheme"
wget https://raw.githubusercontent.com/dylnmc/novum.vim/master/colors/novum.vim 2>/dev/null
mkdir -p ~/.vim/colors
mv novum.vim ~/.vim/colors/novum.vim

mv ~/.vimrc ~/.vimrc.old
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/.vimrcfile 2>/dev/null
mv .vimrcfile ~/.vimrc

sleep 2
echo "Done!"

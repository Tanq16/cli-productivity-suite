#!/usr/bin/zsh
git clone https://github.com/itchyny/lightline.vim ~/.vim/pack/plugins/start/lightline
git clone https://github.com/jiangmiao/auto-pairs.git ~/auto_pairs_vim
mkdir ~/.vim/plugin
cp ~/auto_pairs_vim/plugin/auto-pairs.vim ~/.vim/plugin/
rm -rf ~/auto_pairs_vim

wget https://www.vim.org/scripts/download_script.php\?src_id\=21752
mv 'download_script.php?src_id=21752' ~/.vim/supertab.vmb
vim -c 'so %' -c 'q' ~/.vim/supertab.vmb

echo "set number" >> ~/.vimrc
echo "set tabstop=4" >> ~/.vimrc
echo "set autoindent" >> ~/.vimrc
echo "set expandtab" >> ~/.vimrc
echo "set laststatus=2" >> ~/.vimrc
echo "syntax on" >> ~/.vimrc
echo "set hlsearch" >> ~/.vimrc
echo "nnoremap \\ :noh<return>" >> ~/.vimrc

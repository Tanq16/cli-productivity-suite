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

echo "Installing novum colorscheme"
wget https://raw.githubusercontent.com/dylnmc/novum.vim/master/colors/novum.vim 2>/dev/null
mkdir -p ~/.vim/colors
mv novum.vim ~/.vim/colors/novum.vim

echo "Setting .vimrc"
echo "set mouse=a" >> ~/.vimrc
echo "set number" >> ~/.vimrc
echo "set tabstop=4" >> ~/.vimrc
echo "set autoindent" >> ~/.vimrc
echo "set expandtab" >> ~/.vimrc
echo "set laststatus=2" >> ~/.vimrc
echo "syntax on" >> ~/.vimrc
echo "set hlsearch" >> ~/.vimrc
echo "colorscheme novum" >> ~/.vimrc

echo "let g:lightline = {" >> ~/.vimrc
echo "      \\ 'colorscheme': 'wombat'," >> ~/.vimrc
echo "      \\}" >> ~/.vimrc
echo "set noshowmode" >> ~/.vimrc
echo "nnoremap \\ :noh<return>" >> ~/.vimrc
echo "set rtp+=~/.fzf" >> ~/.vimrc
echo "set cursorline" >> ~/.vimrc
echo "nnoremap <C-m> :tabnext<return>" >> ~/.vimrc
echo "nnoremap <C-n> :tabprev<return>" >> ~/.vimrc
echo "nnoremap ff :FZF ~<return>" >> ~/.vimrc
echo "nnoremap tt :tabnew<return>" >> ~/.vimrc
echo "nnoremap nm :set invnumber<return>" >> ~/.vimrc

echo "if has("autocmd")" >> ~/.vimrc
echo "  au BufReadPost * if line(\"'\\\"\") > 0 && line(\"'\\\"\") <= line(\"$\") | exe \"normal! g\`\\\"\" | endif" >> ~/.vimrc
echo "endif" >> ~/.vimrc

sleep 2
echo "Done!"

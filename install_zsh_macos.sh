#!/bin/sh

echo ""
echo "If you have some other vim config installed, press ⌃+c now and remove that."
echo "After setting up appropriately, start the script again. Sleeping for 20 seconds!"
for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do echo -n '.'; sleep 1; done; echo ""

echo "\n"
echo "Installing several packages; may take a few minutes. Hang tight!\n"

# brew installs
brew install tree tmux jq bat fd neovim lsd wget git 1>/dev/null 2>/dev/null
echo -n '.'

# OMZ and plugins
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1 2>/dev/null
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -ie "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions 2>/dev/null
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting 2>/dev/null
sed -ie "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
echo -n '.'

# NeoVIM setup
rm -rf ~/.vim* 1>/dev/null 2>/dev/null
rm -rf ~/.config/nvim 2>/dev/null
rm -rf ~/.local/share/nvim 2>/dev/null
git clone https://github.com/NvChad/NvChad ~/.config/nvim --depth 1 1>/dev/null 2>/dev/null
sed -ie "s/local input =/local input = N --/" ~/.config/nvim/lua/core/bootstrap.lua
sed -ie "s/dofile(vim.g/vim.cmd([[ set guicursor= ]])\ndofile(vim.g/" ~/.config/nvim/init.lua
nvim --headless -c 'quitall' 1>/dev/null 2>/dev/null
echo -n '.'

# TMUX & FZF
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 1>/dev/null 2>/dev/null
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf 1>/dev/null 2>/dev/null
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 1>/dev/null 2>/dev/null
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf 1>/dev/null 2>/dev/null
~/.fzf/install --all 1>/dev/null 2>/dev/null
echo -n '.'

# RC file setup
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc 2>/dev/null
sed -ie "s/fdfind/fd/" ./add_to_rc
sed -ie "s/batcat/bat/" ./add_to_rc
sed -ie "s/alias bat=/\# alias bat=/" ./add_to_rc
sed -ie "s/alias ip4/alias ip4='ifconfig | grep -oE \"\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\" | grep -v \"127.0.0.1\" | grep -vE \".+\.255$\"' #/" ./add_to_rc
cat ~/.zshrc >> ./temptemp
cat ./add_to_rc >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./add_to_rc ./add_to_rce ./.zshrce 1>/dev/null 2>/dev/null

echo "\n\n"
echo "If you don't see shapes properly after this, make sure to install the font properly (Read the README)"
echo ""
echo "NOTE: After the new shell spawns, quit the terminal app or SSH session for everything to take effect then start again"
echo ""
echo "Starting in 10 seconds!"
for i in 1 2 3 4 5 6 7 8 9 10; do echo -n '.'; sleep 1; done; echo ""
exec zsh -l

su - # become root
usermod -aG sudo tanq # add your user to sudo group
visudo # check that sudo group is allowed for all
# restart the system

sudo apt update -y
# remove bloat
sudo apt remove -y \
	gnome-nibbles gnome-2048 gnome-mahjongg gnome-sudoku \
	gnome-weather gnome-software gnome-software-common gnome-music \
	gnome-mines gnome-chess gnome-calendar gnome-maps rhythmbox \
	gnome-klotski gnome-clocks iagno gnome-robots tali hitori \
	simple-scan aisleriot four-in-a-row lightsoff quadrapassel \
	five-or-more gnome-contacts pegsolitaire gnome-robots \
	gnome-tetravex gnome-taquin swell-foop evolution \
	transmission-* libreoffice-*
sudo apt autoremove && sudo apt autoclean

# initial update cycle
sudo apt upgrade -y && sudo apt dist-upgrade -y

# install nala (better apt frontend) and other init stuff
sudo apt install -y nala openssh-server openssh-client \
	gnome-shell-extension-manager curl wget tilix htop

# disable cups browsed service permanently
sudo systemctl stop cups-browsed.service
sudo systemctl disable cups-browsed.service

# install proper docker (not distro-bundled)
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt remove $pkg; done
sudo nala install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo   "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" |   sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo nala update -y && sudo nala install -y \
	docker-ce docker-ce-cli containerd.io \
	docker-buildx-plugin docker-compose-plugin
sudo groupadd docker
sudo usermod -aG docker $USER

# remove unwanted directories from home
rm -rf Documents/ Music/ Public/ Templates/ Videos/
nano .config/user-dirs.dirs # change unwanted ones to $HOME only

# restart the system

# cli productivity suite
sudo nala install -y \
	tar wget tree tmux jq ninja-build \
	gettext make cmake unzip git file \
	gcc bat fd-find zsh
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh
git clone https://github.com/spaceship-prompt/spaceship-prompt.git ~/.oh-my-zsh/custom/themes/spaceship-prompt --depth=1
ln -s ~/.oh-my-zsh/custom/themes/spaceship-prompt/spaceship.zsh-theme ~/.oh-my-zsh/custom/themes/spaceship.zsh-theme
sed -i "s/robbyrussell/spaceship/" ~/.zshrc
git clone --depth=1 https://github.com/zsh-users/zsh-autosuggestions.git ~/.oh-my-zsh/custom/plugins/zsh-autosuggestions
git clone --depth=1 https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.oh-my-zsh/custom/plugins/zsh-syntax-highlighting
sed -i "s/plugins=/plugins=(git zsh-autosuggestions zsh-syntax-highlighting) #/" ~/.zshrc
rm -rf ~/.vim* ~/.config/nvim ~/.local/share/nvim 
sudo nala remove vim neovim -y
wget https://github.com/neovim/neovim/archive/refs/tags/stable.tar.gz
tar -xvf stable.tar.gz
cd neovim-stable
make CMAKE_BUILD_TYPE=RelWithDebInfo
cd build
cpack -G DEB
sudo nala install ./nvim-linux64.deb
cd ../.. && rm -rf stable.tar.gz neovim-stable
sed -i "s/local input =/local input = \"N\" --/" ~/.config/nvim/lua/core/bootstrap.lua
git clone https://github.com/NvChad/NvChad ~/.config/nvim
sed -i "s/local input =/local input = \"N\" --/" ~/.config/nvim/lua/core/bootstrap.lua
sed -i "s/dofile(vim.g/vim.cmd([[ set guicursor= ]])\ndofile(vim.g/" ~/.config/nvim/init.lua
nvim --headless -c 'quitall'
git clone --depth=1 https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/tmuxconf
mv tmuxconf ~/.tmux.conf
TMUX_PLUGIN_MANAGER_PATH=~/.tmux/plugins ~/.tmux/plugins/tpm/bin/install_plugins 
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install --all
a=$(curl -L -s https://github.com/lsd-rs/lsd/releases/ | grep -oE "tag.+\"" | cut -d '/' -f2 | grep -vE "^[^v]" | cut -d "\"" -f1 | head -n 1)
wget "https://github.com/lsd-rs/lsd/releases/download/$a/lsd-""$a""-x86_64-unknown-linux-gnu.tar.gz"
tar -xvf "lsd-""$a""-x86_64-unknown-linux-gnu.tar.gz"
sudo mv "lsd-""$a""-x86_64-unknown-linux-gnu/lsd" /usr/bin/lsd
rm -rf "lsd-""$a""-x86_64-unknown-linux-gnu.tar.gz" "lsd-""$a""-x86_64-unknown-linux-gnu"
wget https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/add_to_rc
cat ~/.zshrc >> ./temptemp
cat ./add_to_rc >> ./temptemp
cat ./temptemp | grep -vE "^#" | grep -vE "^$" > ~/.zshrc
rm ./temptemp ./add_to_rc
exec zsh -l
# restart the system

# download catpuccin for tilix and jetbrains nerd font
wget https://raw.githubusercontent.com/catppuccin/tilix/main/src/Catppuccin-Mocha.json
mkdir -p ~/.config/tilix/schemes/
mv Catppuccin-Mocha.json ~/.config/tilix/schemes
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/JetBrainsMono.zip
mkdir /tmp/fontest && mv JetBrainsMono.zip /tmp/fontest && cd /tmp/fontest
unzip JetBrainsMono.zip
sudo mkdir /usr/share/fonts/truetype/jetbrains-nerdfont
sudo mv *.ttf /usr/share/fonts/truetype/jetbrains-nerdfont/
sudo fc-cache -fv

# setup custom DNS
sudo apt install resolvconf
sudo systemctl start resolvconf.service
sudo systemctl enable resolvconf.service
sudo nano /etc/resolvconf/resolv.conf.d/head
sudo resolvconf --enable-updates
sudo resolvconf -u

# Debian Configuration

The following set of commands are an easy way to configure a Debian system for primary use.

First, we become root and add our user to the `sudo` group to grant them administrative through `sudo` &rarr;

```bash
su - # become root
usermod -aG sudo tanq # add your user to sudo group
visudo # check that sudo group is allowed for all
```

After this, ***restart the system***.

Next, these two commands remove some basic debian bloat &rarr;

```bash
sudo apt update -y && sudo apt remove -y \
	gnome-nibbles gnome-2048 gnome-mahjongg gnome-sudoku \
	gnome-weather gnome-software gnome-software-common gnome-music \
	gnome-mines gnome-chess gnome-calendar gnome-maps rhythmbox \
	gnome-klotski gnome-clocks iagno gnome-robots tali hitori \
	simple-scan aisleriot four-in-a-row lightsoff quadrapassel \
	five-or-more gnome-contacts pegsolitaire gnome-robots \
	gnome-tetravex gnome-taquin swell-foop evolution \
	transmission-* libreoffice-*
```

```bash
sudo apt autoremove -y && sudo apt autoclean
```

Now trigger an update cycle &rarr;

```bash
sudo apt update -y && sudo apt upgrade -y && sudo apt dist-upgrade -y
```

Now, install a couple tools including `nala` (fast frontend for `apt`) &rarr;

```bash
sudo apt install -y nala openssh-server openssh-client wget \
	gnome-shell-extension-manager curl htop wl-clipboard
```

Disable cups-browsed service permanently to remove the wait during a restart &rarr;

```bash
sudo systemctl stop cups-browsed.service
sudo systemctl disable cups-browsed.service
```

Install docker properly (not the distro bundled version) - copy each command one by one &rarr;

```bash
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt remove $pkg; done
sudo nala install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo   "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" |   sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo nala update && sudo nala install \
	docker-ce docker-ce-cli containerd.io \
	docker-buildx-plugin docker-compose-plugin
sudo groupadd docker
sudo usermod -aG docker $USER
```

Download and install the JetBrainsMono Nerd Font &rarr;

```bash
wget https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/JetBrainsMono.zip && \
mkdir /tmp/fontest && mv JetBrainsMono.zip /tmp/fontest && cd /tmp/fontest && \
unzip JetBrainsMono.zip && \
sudo mkdir /usr/share/fonts/truetype/jetbrains-nerdfont && \
sudo mv *.ttf /usr/share/fonts/truetype/jetbrains-nerdfont/ && \
sudo fc-cache -fv
```

Now some QoL updates &rarr;

```bash
# remove unwanted directories from home
rm -rf Documents/ Music/ Public/ Templates/ Videos/
nano .config/user-dirs.dirs # change unwanted ones to $HOME only
```

After these steps, ***restart the system***. This can be done using `sudo reboot`.

Now, we install the cli productivity suite &rarr;

```bash
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
# install with --- sh install.sh
# cleanup with --- rm install.sh
```

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh_linux.sh)"
```

Now, restart the shell completely. Next, download the `kitty` terminal &rarr;

```bash
curl -L https://sw.kovidgoyal.net/kitty/installer.sh | sh /dev/stdin
```

```bash
ln -sf ~/.local/kitty.app/bin/kitty ~/.local/kitty.app/bin/kitten ~/.local/bin/
cp ~/.local/kitty.app/share/applications/kitty.desktop ~/.local/share/applications/
cp ~/.local/kitty.app/share/applications/kitty-open.desktop ~/.local/share/applications/
sed -i "s|Icon=kitty|Icon=/home/$USER/.local/kitty.app/share/icons/hicolor/256x256/apps/kitty.png|g" ~/.local/share/applications/kitty*.desktop
sed -i "s|Exec=kitty|Exec=/home/$USER/.local/kitty.app/bin/kitty|g" ~/.local/share/applications/kitty*.desktop
```

```bash

```

**Bonus**: Setup custom DNS &rarr;

```bash
sudo apt install resolvconf && \
sudo systemctl start resolvconf.service && \
sudo systemctl enable resolvconf.service && \
sudo nano /etc/resolvconf/resolv.conf.d/head
```

Add nameserver such as `nameserver    192.168.55.32` and then run &rarr;

```bash
sudo resolvconf --enable-updates && sudo resolvconf -u
```

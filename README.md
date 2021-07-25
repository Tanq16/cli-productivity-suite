# Command Line Productivity Suite

Use this script to easily install **zsh** and the custom shell experience of **Oh-My-Zsh**. This also installs `bat` as an alternative to `cat` with syntax hughlighting. A bonus script is also present for quick settings of `vim`.

## Requirements

System should have curl installed. `vim` is required for the bonus script. If not, do it like so:
```bash
sudo apt install curl vim
```

## Installation
The first step is to install zsh. Enter your password when prompted to make zsh default (or do it on your own as per your wish). Execute the following on your terminal -
```bash
sudo apt install -y zsh wget git
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh
```

Once done, execute the following to install all other magic and enter your password whenever prompted.
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh.sh)"
```
Or you could clone the repo and run the `install_zsh.sh` script.

Restart your shell (log out and back in on a linux machine or simply close the shell and start again on something like WSL). Execute `p10k configure`. This time, set the options according to the way you want your prompt to look. For any wierd looking prompts, see Post Installation steps below.

## Features
This script installs the following -
1. Zsh.
2. OH-MY-ZSH custom shell with Powerlevel10k theme.
3. Fuzzy finder for awesome productivity.
3. Syntax highlighting for command line.
4. Auto-completion on command line.
5. Color scheme that looks good on dark terminals.

## Post Installation
If the prompt looks funny or as unintended or options for lock, debian logo, etc. were not visible, please change font of the terminal you use to a powerline-font. Best recommendation from me - FiraCode Nerd Font, Ubuntu Mono Powerline Derivative.

To install from repositories on debian systems, use
```bash
sudo apt install fonts-powerline
```

After this reconfigure the powerlevel10k using -
```bash
p10k configure
```

If you are installing on windows, you could just install the **FiraCode Nerd Font Mono** (recommended) font from [here](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip) or use the the repository to install all fonts as per the instructions there. Download the folder and install Ubuntu Mono Derivative for powerline and set it as your font in whatever terminal you use. Then restart the terminal app. For windows, right click the font ttf and install. The settings json file must be edited for setting this as the default font in windows terminal.

The fuzzy search is another awesome feature to have. Read the specifics at the github page [here](https://github.com/junegunn/fzf) and learn about the features a bit more [here](https://medium.com/better-programming/boost-your-command-line-productivity-with-fuzzy-finder-985aa162ba5d).

## Bonus: Basic Vim quick changes
To use the basic modern editor features in vim, use the given script to install the features as follows -
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/vim_improve.sh)"
```
This installs the supertab (tab => autocomplete), auto-pair brackets, lightline plugin and sets numbering, expandtab, tab=4 spaces and autoindent.

## Bonus: Installation of Bat
Bat is an alternative of the `cat` command. It uses syntax highlighting and works similar to cat. This is installed by default. Bat also enables pager by default i.e., the command is automatically piped to `less` when the output is larger than a threshold. This is disabled in the installion by using -
```bash
export BAT_PAGER=''
```
during the installation. You could re-enable the same by deleting the line in your `.zshrc`. Or you could also just paste the following -
```bash
sed -i "s/export BAT_PAGER=''//" ~/.zshrc
```

# AIO Container
The given `Dockerfile` is for an ubuntu that has essentially all tools above installed and ssh enabled for anyone to reach it via ssh and work on it. The `p10k.zsh` file has to be inside the same directory as the `Dockerfile` as it copies it and prevents the configuration wizard of zsh from running when you run the docker. If the wizard is still needed for customization, then run `p10k configure` and replcae the contents of the `p10k.zsh` file in the host with those of the `~/.p10k.zsh` file inside the docker.

After this, it is possible to ssh into the docker with the root user and password `docker`. To build the docker use this command in the folder that contains the files -
```bash
docker build -t aio_docker .
```
Thereafter, run the following command to run the docker -
```bash
docker run --name="aio_docker_instance" --rm -p 50232:22 -it aio_docker zsh -c "service ssh start; zsh"
```
To ssh into the docker without adding it to the hosts file, use the following command -
```bash
ssh -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null" root@localhost -p 50232
```
More parameters for specifying shared volumes or other ports can be added to the run command.
It is also useful to have buildkit enabled. This can be done by using the follwing command or adding it to the shell rc file -
```bash
export DOCKER_BUILDKIT=1
```
Disabling can be done by making the above 0.

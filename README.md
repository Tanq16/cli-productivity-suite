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

**FiraCode Nerd Font Mono** (recommended) font can be downloaded from [here](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip).

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

# Security focussed docker image

This docker image is intended to be used for cybersecurity related operations. It's effectively a combination of all the good tools required for basic pentesting. This includes the command line enhancements listed above. The intended way to use it is to run the docker and then ssh into the instance via VS code and a terminal application, use tmux and work.

This is specific to x86 machines and the best way to get access is to pull the image by using -
```bash
docker pull tanq16/sec_docker:main
```
Then, it can be run by using -
```bash
docker run --name="sec_docker" -v ~/go/:/root/go/ --rm -p 50022:22 -it tanq16/sec_docker:main zsh -c "service ssh start; zsh"
```

After this, it is possible to ssh into the docker with the `root` user and password `docker`.

The above command mounts the local volume for go code. Similarly, the only port mapped to localhost is 50022, for ssh. Jupyter-Lab is also installed on the docker and the port 8888 can be mapped to that of localhost:58888. The general norm for mapping ports for this docker should be 50000+port for consistency and interoperability.

## Notable installations in the docker

nmap, ncat & ncrack
ltrace & strace
gobuster, nikto & dirb
netdiscover & wireshark (tshark mainly, because its cli)
hydra, fcrackzip & john the ripper
pwndbg
metasploit-framework & searchsploit (with exploit-database)
jupyter-lab & golang
seclists & rockyou.txt

## Build on your own

The repository includes the required files to build the image. If the exact tools are not required and extra tools must be installed or replaced, the given `Dockerfile` should be edited. 

The `p10k.zsh` file has to be inside the same directory as the `Dockerfile` as it copies it and prevents the configuration wizard of zsh from running when you run the docker. If the wizard is still needed for customization, then run `p10k configure` inside the docker and replace the contents of the `p10k.zsh` file in the host with those of the `~/.p10k.zsh` file inside the docker.

To build the docker use this command in the folder that contains the files -
```bash
docker build -t aio_docker .
```
Thereafter, run the following command to run the docker -
```bash
docker run --name="aio_docker_instance" --rm -p 50022:22 -it aio_docker zsh -c "service ssh start; zsh"
```

## Bonus information

To ssh into the docker without adding it to the hosts file, use the following command -
```bash
ssh -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null" root@localhost -p 50232
```
or add the following alias to the rc file for your default shell -
```bash
alias sshide='ssh -o "StrictHostKeyChecking=no" -o "UserKnownHostsFile=/dev/null"'
```

It is also useful to have buildkit enabled. This can be done by using the following command or adding it to the shell rc file -
```bash
export DOCKER_BUILDKIT=1
```
Disabling can be done by making the above 0.

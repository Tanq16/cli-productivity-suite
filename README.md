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

# Docker Images

The two docker images in this repository are intended to be used for cybersecurity related operations and for python/go development. The security docker is effectively a combination of all the good tools required for basic pentesting. Both the images include the command line enhancements listed above. The intended way to use it is to run the docker and then ssh into the instance via VS code and a terminal application, use tmux and work.

This is specific to x86 machines and the best way to get access is to pull the image by using -
```bash
docker pull tanq16/sec_docker:main # For the security docker image/
docker pull tanq16/sec_docker:dev # For the development docker image
```
Then, it can be run by using -
```bash
# The dev docker image - mount the programming folders from host and map the ssh port
docker run --name="sec_docker" -v ~/go_programs/:/root/go/src -v ~/python_programs/:/root/python/ --rm -p 50022:22 -it tanq16/sec_docker:dev zsh -c "service ssh start; zsh"

# The security docker image - map the ssh port and the jupyterlab port
docker run --name="sec_docker" --rm -p 58080:8080 -p 50022:22 -it tanq16/sec_docker:main zsh -c "service ssh start; zsh"
```

The security image also has the development instructions in its `Dockerfile`, so the volumes can be mounted there as well. On connecting the VS code via the remote ssh extension to the docker image, the python package and the go package

The `service ssh start` section of the command to be executed is needed to enable ssh access. Direct loading of the shell interferes with the oh-my-zsh themes and not all things are loaded. Therefore, the docker image should be run either in background or as stated above to signify a control shell and then use ssh and tmux to simulate work environment. After this, it is possible to ssh into the docker with the `root` user and password `docker`.

The general norm for mapping ports for these images is 50000+port for consistency and interoperability.

## Notable installations in the security docker

* nmap, ncat & ncrack
* ltrace & strace
* gobuster, nikto & dirb
* netdiscover & wireshark (tshark mainly, because its cli)
* hydra, fcrackzip & john the ripper
* pwndbg
* metasploit-framework & searchsploit (with exploit-database)
* jupyter-lab & golang
* seclists & rockyou.txt

## Build on your own

The repository includes the required files to build both the images. If the existing tools are not required or extra tools must be installed or replaced, the given `Dockerfile`s in the respective directories should be edited. 

The `p10k.zsh` file for each directory must be inside the same directory as the `Dockerfile` of the respective docker image, as the vuild process copies it and prevents the configuration wizard of oh-my-zsh from running when accessing the shell of the docker image. If the wizard is still needed for customization, then run `p10k configure` inside the docker and replace the contents of the `p10k.zsh` file in the host with those of the `~/.p10k.zsh` file inside the directory for the required docker image.

To build the docker use this command in the required folder that contains the files -
```bash
docker build -t dev_docker .
```
Thereafter, run the following command to execute the shell within the image -
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

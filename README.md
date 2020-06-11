# Fast Installation of Oh-My-Zsh

Use this script to easily install **zsh** and the custom shell experience of **Oh-My-Zsh**.

## Requirement

System should have curl installed.
```bash
sudo apt install curl
```

## Installation
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/oh-my-zsh-speed-installation/master/install_zsh.sh)"
```
The script also asks if you want to install [powerlevel10k](https://github.com/romkatv/powerlevel10k).
Select option `1` to install OhMyZsh (Recommended and usually great for an everyday user).
Select `2` if you need more customization along with additional features (which you need to install on your own in addition to autocomplete and suggestions that this script installs).

## Features
This script installs the following -
1. Zsh
2. OH-MY-ZSH custom shell (also gives an option to install Powerlevel10k)
3. Syntax highlighting for command line
4. Auto-completion on command lilne

## Post Installation
If the prompt looks funny or as unintended, please change font of the terminal you use to a powerline-font. Best recommendation from me - Ubuntu Mono Powerline Derivative. 
To install from repositories, use
```bash
sudo apt install fonts-powerline
```

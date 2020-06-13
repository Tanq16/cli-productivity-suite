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
## Features
This script installs the following -
1. Zsh.
2. OH-MY-ZSH custom shell with Powerlevel10k.
3. Syntax highlighting for command line
4. Auto-completion on command line

## Post Installation
If the prompt looks funny or as unintended, please change font of the terminal you use to a powerline-font. Best recommendation from me - Ubuntu Mono Powerline Derivative. 
To install from repositories, use
```bash
sudo apt install fonts-powerline
```
If you are installing on windows, Install the Ubuntu Mono Derivative as the recommended font from [here](https://github.com/powerline/fonts/blob/master/UbuntuMono/). Download the folder and install Ubuntu Mono Derivative for powerline and set it as your font in whatever terminal you use. Then restart the terminal app.
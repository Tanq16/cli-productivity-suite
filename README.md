# Fast Installation of Oh-My-Zsh

Use this script to easily install **zsh** and the custom shell experience of **Oh-My-Zsh**.

## Requirements

System should have curl installed. If not, do it like so:
```bash
sudo apt install curl
```

## Installation
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/oh-my-zsh-speed-installation/master/install_zsh.sh)"
```
Or you could clone the repo and run the `install_zsh.sh` script.

## Features
This script installs the following -
1. Zsh.
2. OH-MY-ZSH custom shell with Powerlevel10k theme.
3. Syntax highlighting for command line.
4. Auto-completion on command line.
5. Color scheme that looks good on dark terminals.

## Post Installation
If the prompt looks funny or as unintended, please change font of the terminal you use to a powerline-font. Best recommendation from me - Ubuntu Mono Powerline Derivative.

To install from repositories on debian systems, use
```bash
sudo apt install fonts-powerline
```

If you are installing on windows, you could just install the **FuraCode Nerd Font Mono** (recommended) font from [here](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip) or use the the repository to install all fonts as per the instructions there. Download the folder and install Ubuntu Mono Derivative for powerline and set it as your font in whatever terminal you use. Then restart the terminal app. For windows, right click the font ttf and install. The settings json file must be edited for setting this as the default font in windows terminal.

## Bonus: Basic Vim quick changes
To use the basic modern editor features in vim, use the given script to install the features as follows -
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/oh-my-zsh-speed-installation/master/vim_improve.sh)"
```
This installs the supertab (tab => autocomplete), auto-pair brackets, lightline plugin and sets numbering, expandtab, tab=4 spaces and autoindent.

# Command Line Productivity Suite

* [Introduction](#introduction)
* [Installation](#installation)
* [Post Installation](#post-installation)
* [Bonus (Helpful Tips)](#bonus-helpful-tips)

## Introduction

Use this repo to easily install a custom, cool and funky shell experience along with an awesome `neovim` and `tmux` installation. Before anything else though, install the [Catppuccin](https://catppuccin-website.vercel.app/) theme (with the `Mocha` configuration) for your terminal to get the most seamless experience with `tmux` and `neovim`.

Also, install a "nerd" font for your terminal emulator. My recommendation is [JetBrains Mono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.0.2/JetBrainsMono.zip). Don't forget to read the "Post Installation" section.

The scripts have been tested on x86-64 variants of Debian and Ubuntu. The MacOS script works for both Intel and M-series chips. MacOS does need `brew` to be installed however.

## Installation

Get started by installing the minimum set of tools &rarr;

```bash
sudo apt install git zsh wget curl
```

If you're on MacOS, use this &rarr;

```bash
brew install git zsh wget curl
```

Next, install oh my zsh as follows &rarr;

```bash
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh # cleanup
```

Then, execute the following to install all other magic and enter the password whenever (if) prompted.

For Linux, use &rarr;

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh_linux.sh)"
```

For MacOS, use &rarr;

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh_macos.sh)"
```

Finally, close the shell ***completely*** (close the terminal app or end the SSH session) and start a new instance.

## Post Installation

- `tmux` is installed by default with the above script. Use `tt` to launch a default session.
- `bat`, an alternative of the `cat` command with colored output is also installed by default. 
- `nvim` is installed with NvChad configuration, but `nvim` doesn't allow setting a theme when running headless, so use `<space>+th` to launch the theme selector inside `nvim`, type and select `catppuccin` to match everything up.

`bat` (maintainer-default) has pager enabled, which is disabled by the installation script using `export BAT_PAGER=''` within the rc-file. This can be re-enabled by deleting the line in `.zshrc`.

## Bonus (Helpful Tips)

A handy shortcut in `tmux` added by the above scripts is `Alt + \` to split into two vertical panes and `Alt + Shift + \` to split into two horizontal panes. Focus can be navigated among the split panes by using `Shift + <arrow keys>`.

Pasting on modified zsh shell can be slow due to magic functions that `oh-my-zsh` installs. A quick fix is to comment those functions in `~/.oh-my.zsh/lib/misc.zsh`. This can also be easily done via the `sed`. The following can be pasted in a file and run as `bash <file>` or `zsh <file>` or as an executable after chmoding the file.

> Replace `-i` flags with `-i'' -e` flags for MacOS

```bash
#!/bin/zsh
sed -i "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
```

If something goes wrong or you see an error during installation, you can remove everything with the following command from the home directory and start from scratch again &rarr;

```bash
rm -rf .oh-my-zsh .fzf .fzf.zsh .tmux .tmux.conf .tmux-themepack .vim* .config/nvim .local/share/nvim .zshrc
```

Lastly, the commands list for setting up a base Debian machine from scratch for optimum readiness is given in the `Configure Debian from Scratch` markdown file.

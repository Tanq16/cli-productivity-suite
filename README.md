# Command Line Productivity Suite

* [Introduction](#introduction)
* [Installation](#installation)
* [Post Installation](#post-installation)
* [Bonus Tips](#bonus-tips)

## Introduction

Use this repo to easily install a custom, cool and funky shell experience along with an awesome `neovim` and `tmux` config. This repo also carries a config file for the recommended `Kitty` terminal.

First, make sure you have the [Catppuccin](https://catppuccin-website.vercel.app/) theme (with the `Mocha` configuration) for your terminal. Then, install a `nerd` font for your terminal. My recommendation is [JetBrains Mono Nerd Font](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/JetBrainsMono.zip).

This suite has been tested on and works for x86-64 Debian, x86-64 Ubuntu, x86-64 MacOS, and Apple Silicon MacOS. (PS: MacOS should alreadt have `brew` installed.)

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
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null && sh install.sh
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
- `nvim` is installed with NvChad configuration, but `nvim` needs a couple other small steps to get up and running smoothly &rarr;
    - First, run the Vim command `:MasonInstallAll`
    - Next, run Vim command `:Lazy sync` and once again run `:MasonInstallAll`
    - NeoVim doesn't allow setting a theme while headless, so use `<space>+th` and select `catppuccin` to match everything
- Run cleanup as follows &rarr; `rm install.sh && rm -rf fzf`

PS: `bat` (maintainer-default) has pager enabled, which is disabled by the installation script using `export BAT_PAGER=''` within the rc-file. This can be re-enabled by deleting that line in `.zshrc`.

## Bonus Tips

A handy shortcut in `tmux` added by the above scripts is `Alt + \` to split into two vertical panes and `Alt + Shift + \` to split into two horizontal panes. Focus can be navigated among the split panes by using `Shift + <arrow keys>`.

Pasting on modified zsh shell can be slow due to magic functions that `oh-my-zsh` installs. A quick fix is to comment those functions in `~/.oh-my.zsh/lib/misc.zsh`. This is already done by this suite. If you need to re-enable this, simply uninstall the suite, restart your terminal and comment the necessary `sed` lines in the scripts before re-installing.

If something goes wrong or you you want to re-install, you can remove everything with the following command from the home directory and start from scratch again (i.e., the installation section, theme and font will be already configured) &rarr;

```bash
rm -rf .oh-my-zsh .fzf .fzf.zsh .tmux .tmux.conf .tmux-themepack .vimrc .viminfo .vim .config/nvim .local/share/nvim .zshrc
```

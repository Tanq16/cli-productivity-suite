# Command Line Productivity Suite

Use this script to easily install **zsh** and the custom shell experience of **Oh-My-Zsh**. This also installs `bat` as an alternative to `cat` with syntax highlighting. As a bonus, SpaceVim is also installed for a flashy `vim` experience. A nice Tmux configuration is also installed.

Vim and Tmux are both configured to use the [Nord](https://www.nordtheme.com/) theme by default, so it's better to install that same theme to your terminal emulator for a seamless experience. Nord can be installed for all general purpose terminal emulators like iTerm2, Windows Terminal, etc. Or, modify the configurations for your preferred theme.

Don't forget to read "Post Installation" section.

## Requirements

The system should have `git`, `zsh`, `wget`, `curl` and `vim` installed. If not, do it like so &rarr;

```bash
sudo apt install git zsh wget curl vim # Or `brew install git zsh wget curl vim` for MacOS
```

## Installation

The first step is to install `zsh`. Enter the password when prompted to make `zsh` default (or do it later as per necessity). Execute the following on your terminal &rarr;

```bash
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh # cleanup
```

Once done, execute the following to install all other magic and enter the password whenever prompted.

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh.sh)"
```

Or, clone the repo and run the `install_zsh.sh` script. Finally, close the shell and start a new instance. If something goes wrong, delete everything using &rarr;

```bash
rm -rf .oh-my-zsh .fzf .fzf.zsh .tmux .tmux.conf .tmux-themepack .vim .vimrc .SpaceVim .SpaceVim.d .zshrc
# The old ZSH rc file may be postfixed with `.pre-oh-my-zsh` or something similar.
```

## Features

This script installs the following &rarr;

1. Oh-My-Zsh custom shell with [spaceship-prompt](https://spaceship-prompt.sh/) theme
2. Fuzzy finder (`fzf`) for awesome productivity
3. Syntax highlighting for command line
3. Auto-completion on command line
4. Tmux with mouse and other quality of life improvements
5. SpaceVim for a flashy vim experience
6. Nord theme for tmux and vim

## Post Installation

Highly recommended font for this suite &rarr; **FiraCode Nerd Font Mono**, which can be downloaded from [here](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip).

The fuzzy search by `fzf` is another awesome feature to have and is installed as a part of the scripts above. Read the specifics at the github page [here](https://github.com/junegunn/fzf) and learn about the features a bit more [here](https://medium.com/better-programming/boost-your-command-line-productivity-with-fuzzy-finder-985aa162ba5d).

Tmux is another addition in the scripts above. To make full use of the plugins installed, the first run must be used to install those plugins. This can be done by using the aliases set in the rc-file. Use `tt` to launch a default session. After that, use `Ctrl + b` followed by `Shift + i` to install the plugins. This takes 3-4 seconds after which the changes take effect and the awesome tmux can be used. A very handy shortcut in the configuration file added by the above scripts is `Alt + \` to split into two vertical panes and `Alt + Shift + \` to split into two horizontal panes. Focus can be navigated among the split panes by using `Shift + <arrow keys>`.

**Tmux on iTerm** can be very useful if certain settings are enabled to allow copy to Mac clipboard upon selection of content within the tmux terminal draw. Tmux plugins installed as a part of the scripts above enable mouse support.

## Bonus: Bat Config

Bat is an alternative of the `cat` command. It uses syntax highlighting and works similar to cat. Bat also enables pager by default i.e., the command is automatically piped to `less` when the output is larger than a threshold. This is disabled with the install script by using &rarr;

```bash
export BAT_PAGER=''
```

During the installation. You could re-enable the same by deleting the line in your `.zshrc` or run the following command &rarr;

```bash
sed -i "s/export BAT_PAGER=''//" ~/.zshrc
# `sed -ie "s/export BAT_PAGER=''//" ~/.zshrc` for MacOS
```

## Bonus: Fix slow Oh-My-Zsh Paste

Pasting on modified zsh shell can be slow due to magic functions that oh-my-zsh installs. A quick fix is to comment those functions in `~/.oh-my.zsh/lib/misc.zsh`. This can also be easily done via the `sed`. The following can be pasted in a file and run as `bash file` or `zsh file` or as an executable after chmoding the file.

```bash
#!/bin/zsh
# Replace `-i` flags with `-ie` flags for MacOS
sed -i "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
```

# Command Line Productivity Suite

* [Introduction](#introduction)
* [Installation](#installation)
* [Post Installation](#post-installation)
* [Bonus (Helpful Tips)](#bonus-helpful-tips)

## Introduction

Use this to easily install a custom and funky shell experience along with a cool `vim` and `tmux` installation. Before anything else though, install the [Nord](https://www.nordtheme.com/) theme for the most seamless experience with `tmux` and `vim`. You could also change configurations for other preferred themes.

Also install a "nerd" font for you terminal emulator, because ligatures are cool. Recommended one is [Fira Code Nerd Font Mono](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip).

Don't forget to read the "Post Installation" section.

<details>
<summary>Everything that the script installs is in this expandable block.</summary>

1. Oh-My-Zsh custom shell with [spaceship-prompt](https://spaceship-prompt.sh/) theme
2. Fuzzy finder (`fzf`) for awesome productivity
3. Syntax highlighting for command line
3. Auto-completion on command line
4. Tmux with mouse and other quality of life improvements
5. SpaceVim for a flashy vim experience
6. Nord theme for tmux and vim

</details>

## Installation

The system should have `git`, `zsh`, `wget`, `curl` and `vim` installed. If not, do it like so &rarr;

```bash
sudo apt install git zsh wget curl vim # Or `brew install git zsh wget curl vim` for MacOS
```

Next, install oh my zsh as follows &rarr;

```bash
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh # cleanup
```

Then, execute the following to install all other magic and enter the password whenever (if) prompted &rarr;

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh.sh)"
```

Finally, close the shell completely and start a new instance.

If something goes wrong or you see an error, remove everything with the following &rarr;

```bash
rm -rf .oh-my-zsh .fzf .fzf.zsh .tmux .tmux.conf .tmux-themepack .vim .vimrc .SpaceVim .SpaceVim.d .zshrc
# The old ZSH rc file may be postfixed with `.pre-oh-my-zsh` or something similar.
```

## Post Installation

`tmux` is installed by default with the above script. To make full use of the plugins installed, the first run must be used to install those plugins. This can be done by using the aliases set in the rc-file as follows &rarr;

* Use `tt` to launch a default session
* Use `Ctrl + b` followed by `Shift + i` to install the plugins
* This takes 3-4 seconds after which the changes take effect and the awesome tmux can be used

`bat`, an alternative of the `cat` command with colored output is also installed by default. `bat` (maintainer-default) has pager enabled, which is disabled by the above script using &rarr;

```bash
export BAT_PAGER=''
```

within the rc-file. This can be re-enabled by deleting the line in `.zshrc` or running the following command &rarr;

```bash
sed -i "s/export BAT_PAGER=''//" ~/.zshrc
# `sed -ie "s/export BAT_PAGER=''//" ~/.zshrc`     (for MacOS) but this will save .zshrce as backup
```

## Bonus (Helpful Tips)

A handy shortcut in `tmux` added by the above scripts is `Alt + \` to split into two vertical panes and `Alt + Shift + \` to split into two horizontal panes. Focus can be navigated among the split panes by using `Shift + <arrow keys>`.

**Tmux on iTerm** can be very useful if certain settings are enabled to allow copy to Mac clipboard upon selection of content within the `tmux` terminal draw. `tmux` plugins installed as a part of the scripts above also enable mouse support.

The fuzzy search `fzf` is another awesome feature to have and is installed as a part of the scripts above. Read the specifics at the maintainer [github page](https://github.com/junegunn/fzf) and learn about the features a bit more [here](https://medium.com/better-programming/boost-your-command-line-productivity-with-fuzzy-finder-985aa162ba5d).

Pasting on modified zsh shell can be slow due to magic functions that `oh-my-zsh` installs. A quick fix is to comment those functions in `~/.oh-my.zsh/lib/misc.zsh`. This can also be easily done via the `sed`. The following can be pasted in a file and run as `bash <file>` or `zsh <file>` or as an executable after chmoding the file.

```bash
#!/bin/zsh
# Replace `-i` flags with `-ie` flags for MacOS
sed -i "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
```

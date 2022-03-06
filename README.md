# Command Line Productivity Suite

Use this script to easily install **zsh** and the custom shell experience of **Oh-My-Zsh**. This also installs `bat` as an alternative to `cat` with syntax highlighting. A bonus script is also present for quick settings of `vim`.

## Requirements

The system should have `curl` installed. `vim` is required for the bonus script. If not, do it like so&rarr;
```bash
sudo apt install curl vim # Or `brew install curl vim` for MacOS
```

## Installation

The first step is to install `zsh`. Enter the password when prompted to make `zsh` default (or do it later as per necessity). Execute the following on your terminal &rarr;
```bash
sudo apt install -y zsh wget git # Or `brew install wget git` for MacOS
wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh 2>/dev/null
sh install.sh
rm install.sh
```

Once done, execute the following to install all other magic and enter the password whenever prompted.
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/install_zsh.sh)"
```
Or, clone the repo and run the `install_zsh.sh` script.

Restart the shell (log out and back in on a linux machine or close the shell and start again on WSL). Execute `p10k configure` after logging in (preferably after installation of the font). This time, set the options according to the way the prompt must look. For any wierd looking prompts, see Post Installation steps below.

## Features

This script installs the following -
1. Zsh
2. OH-MY-ZSH custom shell with Powerlevel10k theme
3. Fuzzy finder for awesome productivity
3. Syntax highlighting for command line
4. Auto-completion on command line
5. Tmux with mouse and other quality of life improvements

## Post Installation

If the prompt looks funny or has an unintended look or the options for lock, debian logo, etc. were not visible, change the font of the terminal you use to a powerline-font or something that has ligatures and icon support. Highest recommendation from the author &rarr; FiraCode Nerd Font, Ubuntu Mono Powerline Derivative.

To install powerline fonts from repositories on debian systems, use &rarr;
```bash
sudo apt install fonts-powerline
```

After this reconfigure the powerlevel10k using &rarr;
```bash
p10k configure
```

**FiraCode Nerd Font Mono** (recommended) font can be downloaded from [here](https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/FiraCode.zip).

The fuzzy search by `fzf` is another awesome feature to have and is installed as a part of the scripts above. Read the specifics at the github page [here](https://github.com/junegunn/fzf) and learn about the features a bit more [here](https://medium.com/better-programming/boost-your-command-line-productivity-with-fuzzy-finder-985aa162ba5d).

Tmux is another addition in the scripts above. To make full use of the plugins installed, the first run must be used to install those plugins. This can be done by using the aliases set in the rc-file. Use `tt` to launch a default session. After that, use `Ctrl + b` followed by `Shift + i` to install the plugins. This takes 3-4 seconds after which the changes take effect and the awesome tmux can be used. A very handy shortcut in the configuration file added by the above scripts is `Alt + \` to split into two vertical panes and `Alt + Shift + \` to split into two horizontal panes. Focus can be navigated among the split panes by using `Shift + <arrow keys>`.

**Tmux on iTerm** can be very useful if certain settings are enabled to allow copy to Mac clipboard upon selection of content within the tmux terminal draw. Tmux plugins installed as a part of the scripts above enable mouse support.

## Bonus: Basic Vim quick changes

To use the basic modern editor features in vim, use the given script to install the features as follows &rarr;
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/Tanq16/cli-productivity-suite/master/vim_improve.sh)"
```
This installs the supertab (tab => autocomplete), auto-pair brackets, lightline plugin and sets numbering, expandtab, tab=4 spaces, autoindent and some shortcuts.

## Bonus: Installation of Bat

Bat is an alternative of the `cat` command. It uses syntax highlighting and works similar to cat. This is installed by default. Bat also enables pager by default i.e., the command is automatically piped to `less` when the output is larger than a threshold. This is disabled in the installion by using &rarr;
```bash
export BAT_PAGER=''
```
during the installation. You could re-enable the same by deleting the line in your `.zshrc`. Or you could also just paste the following &rarr;
```bash
sed -i "s/export BAT_PAGER=''//" ~/.zshrc
```

## Bonus: Fix slow Oh-My-Zsh Paste

Pasting on modified zsh shells can be slow ue to magic functions that oh-my-zsh installs. A quick fix for this is to comment those functions in the file at location `$HOME/.oh-my.zsh/lib/misc.zsh`. This can also be easily done via the `sed`. The following can be pasted in a file and run as `bash file` or `zsh file` or as an executable after chmoding the file.

```bash
#!/bin/zsh
sed -i "s/autoload -Uz bracketed-paste-magic/#autoload -Uz bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N bracketed-paste bracketed-paste-magic/#zle -N bracketed-paste bracketed-paste-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/autoload -Uz url-quote-magic/#autoload -Uz url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
sed -i "s/zle -N self-insert url-quote-magic/#zle -N self-insert url-quote-magic/" ~/.oh-my-zsh/lib/misc.zsh
```

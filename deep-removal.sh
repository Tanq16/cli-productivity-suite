#!/usr/bin/env bash
# Deep removal of CPS + CPS-installed brew packages + Oh My Zsh.
# Homebrew itself is preserved. ~/.zsh_history is preserved.
#
# After running this, you can start fresh by:
#   1. Reinstalling Oh My Zsh (Homebrew stays)
#   2. Reinstalling the cps binary
#   3. Running `cps init`

echo "CPS deep removal"
echo ""
echo "Will remove:"
echo "  - cps binary (~/.local/bin/cps)"
echo "  - CPS directories (~/shell, ~/.tmux, ~/.config/nvim, ~/.config/cps, ~/.local/share/nvim)"
echo "  - CPS-deployed configs (.zshrc, .tmux.conf, .aerospace.toml, kitty configs)"
echo "  - Oh My Zsh (~/.oh-my-zsh)"
echo "  - Brew packages installed by CPS (neovim, cloud CLIs, core/dev/network/media tools)"
echo ""
echo "Will preserve:"
echo "  - ~/.zsh_history"
echo "  - Homebrew itself"
echo "  - System tools (git, curl, zsh from apt or macOS built-in)"
echo ""
read -rp "Continue? [y/N] " ans
case "$ans" in
  [yY]|[yY][eE][sS]) ;;
  *) echo "aborted"; exit 0 ;;
esac

echo ""
echo "==> removing cps binary"
rm -f "$HOME/.local/bin/cps"

echo "==> removing CPS-managed directories"
rm -rf "$HOME/shell"
rm -rf "$HOME/.tmux"
rm -rf "$HOME/.config/nvim"
rm -rf "$HOME/.config/cps"
rm -rf "$HOME/.local/share/nvim"

echo "==> removing CPS-deployed configs"
rm -f "$HOME/.tmux.conf"
rm -f "$HOME/.zshrc"
rm -f "$HOME/.aerospace.toml"
rm -f "$HOME/.config/kitty/kitty.conf"
rm -f "$HOME/.config/kitty/current-theme.conf"

if command -v brew >/dev/null 2>&1; then
  echo "==> uninstalling CPS-installed brew formulas"
  brew uninstall --ignore-dependencies \
    wget zip unzip file tmux htop neovim \
    cmake gcc make ninja gettext \
    nmap openssl ffmpeg \
    awscli azure-cli \
    2>/dev/null || true

  echo "==> uninstalling CPS-installed brew casks"
  brew uninstall --cask --force gcloud-cli 2>/dev/null || true
  brew uninstall --cask --force nikitabobko/tap/aerospace 2>/dev/null || true
else
  echo "==> brew not found, skipping brew package uninstall"
fi

echo "==> removing Oh My Zsh"
rm -rf "$HOME/.oh-my-zsh"

echo ""
echo "done."
echo "preserved: ~/.zsh_history, Homebrew itself"
echo ""
echo "to start fresh:"
echo "  1. reinstall Oh My Zsh:"
echo "     sh -c \"\$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\""
echo "  2. reinstall cps:"
echo "     see https://github.com/tanq16/cli-productivity-suite#install"
echo "  3. cps init"

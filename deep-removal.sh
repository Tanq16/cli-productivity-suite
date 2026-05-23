#!/usr/bin/env bash
# Deep removal of CPS + CPS-installed brew packages.
# Homebrew itself is preserved. ~/.zsh_history is preserved.
#
# After running this, you can start fresh by:
#   1. Reinstalling the cps binary (Homebrew stays)
#   2. Running `cps init`

echo "CPS deep removal"
echo ""
echo "Will remove:"
echo "  - cps binary (~/.local/bin/cps)"
echo "  - CPS directories (~/shell — includes go-sdk, java-sdk, rust, fnm, py-default,"
echo "    uv tools, all installed binaries; plus ~/.tmux, ~/.config/nvim, ~/.config/cps)"
echo "  - Neovim caches/state (~/.local/share/nvim, ~/.local/state/nvim, ~/.cache/nvim)"
echo "  - Legacy pre-v1.3 paths (~/.nvm, ~/google-cloud-sdk, ~/nuclei-templates)"
echo "  - Legacy runtime caches (~/.local/share/uv, ~/.bun, ~/.npm, go-build cache)"
echo "  - CPS-deployed configs (.zshrc, .tmux.conf, .aerospace.toml, kitty configs, starship.toml)"
echo "  - Brew packages installed by CPS (neovim, nmap, openssl, ffmpeg, aws-cli,"
echo "    azure-cli, gcloud-cli cask)"
echo ""
echo "Will preserve:"
echo "  - ~/.zsh_history"
echo "  - Homebrew itself"
echo "  - System tools (git, curl, zsh from apt or macOS built-in)"
echo "  - Broadly-useful brew formulas (wget, zip, unzip, file, tmux, htop)"
echo "  - Linux dev toolchain (cmake, gcc, make, ninja, gettext)"
echo "  - Aerospace (macOS tiling WM cask)"
echo "  - Custom-extension state (clean those with: cps extend <pack> --remove)"
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

echo "==> removing Neovim caches and state"
rm -rf "$HOME/.local/share/nvim"
rm -rf "$HOME/.local/state/nvim"
rm -rf "$HOME/.cache/nvim"

echo "==> removing legacy pre-v1.3 install locations"
rm -rf "$HOME/.nvm"                  # superseded by fnm in ~/shell/fnm
rm -rf "$HOME/google-cloud-sdk"      # gcloud, before brew cask
rm -rf "$HOME/nuclei-templates"      # moved to ~/shell/nuclei-templates

echo "==> removing legacy runtime caches outside ~/shell"
rm -rf "$HOME/.local/share/uv"       # uv interpreters, before UV_PYTHON_INSTALL_DIR
rm -rf "$HOME/.bun"                  # bun globals/cache, before BUN_INSTALL
rm -rf "$HOME/.npm"                  # npm cache, before npm_config_cache
rm -rf "$HOME/.cache/go-build"       # Linux: go build cache, before GOCACHE
rm -rf "$HOME/Library/Caches/go-build"  # macOS: go build cache, before GOCACHE

echo "==> removing CPS-deployed configs"
rm -f "$HOME/.tmux.conf"
rm -f "$HOME/.zshrc"
rm -f "$HOME/.zprofile"
rm -f "$HOME/.aerospace.toml"
rm -f "$HOME/.config/kitty/kitty.conf"
rm -f "$HOME/.config/kitty/current-theme.conf"
rm -f "$HOME/.config/starship.toml"

if command -v brew >/dev/null 2>&1; then
  echo "==> uninstalling CPS-installed brew formulas"
  # Kept (broadly useful, not CPS-specific): wget zip unzip file tmux htop
  # Kept (Linux dev-tools group): cmake gcc make ninja gettext
  brew uninstall \
    neovim \
    nmap openssl ffmpeg \
    awscli azure-cli \
    2>/dev/null || true

  echo "==> uninstalling CPS-installed brew casks"
  # Kept (broadly useful, not CPS-specific): nikitabobko/tap/aerospace
  brew uninstall --cask --force gcloud-cli 2>/dev/null || true
else
  echo "==> brew not found, skipping brew package uninstall"
fi

echo "==> removing legacy Oh My Zsh install (if present from pre-v1.x CPS)"
rm -rf "$HOME/.oh-my-zsh"

echo ""
echo "done."
echo "preserved: ~/.zsh_history, Homebrew itself"
echo ""
echo "to start fresh:"
echo "  1. reinstall cps:"
echo "     see https://github.com/tanq16/cli-productivity-suite#install"
echo "  2. cps init"

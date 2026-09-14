#!/usr/bin/env bash
# Deep removal of CPS + CPS-installed brew packages.
# Homebrew itself is preserved. ~/.zsh_history is preserved.
#
# After running this, you can start fresh by:
#   1. Reinstalling the cps binary (Homebrew stays)
#   2. Running `cps prereq` then `cps shell`

echo "CPS deep removal"
echo ""
echo "Will remove:"
echo "  - cps binary (~/.local/bin/cps)"
echo "  - CPS directories (~/shell - includes go-sdk, java-sdk, rust, fnm, py-default,"
echo "    uv tools, all installed binaries; plus ~/.tmux, ~/.config/nvim, ~/.config/cps)"
echo "  - Neovim caches/state (~/.local/share/nvim, ~/.local/state/nvim, ~/.cache/nvim)"
echo "  - App state outside ~/shell (~/.local/share/code-server, ~/.config/code-server,"
echo "    ~/.local/kitty.app on Linux, ~/.cache/uv)"
echo "  - Legacy pre-v1.3 paths (~/.nvm, ~/google-cloud-sdk, ~/nuclei-templates)"
echo "  - Legacy runtime caches (~/.local/share/uv, ~/.bun, ~/.npm, go-build cache)"
echo "  - CPS-deployed configs (.zshrc, .tmux.conf, .aerospace.toml, kitty configs, starship.toml)"
echo "  - Brew packages installed by CPS (nmap, openssl, ffmpeg, imagemagick,"
echo "    aws-cli, azure-cli, gcloud-cli cask, plus legacy brew neovim)"
echo ""
echo "Will preserve:"
echo "  - ~/.zsh_history"
echo "  - Homebrew itself"
echo "  - System tools (git, curl, zsh from apt or macOS built-in)"
echo "  - Prerequisites from \`cps prereq\` (wget, zip, unzip, file, tmux, htop,"
echo "    and on Linux the apt build toolchain)"
echo "  - macos-desktop casks (aerospace, sol, font-jetbrains-mono-nerd-font)"
echo "    - remove them by hand if you want them gone"
echo "  - kitty on macOS (/Applications/kitty.app)"
echo "    - remove it by hand if you want it gone"
echo "  - Package state outside ~/shell (~/.local/share/cursor-agent, ~/.config/antigravity)"
echo "  - Neo4j databases and config (~/.config/neo4j) - delete it by hand if you want them gone"
echo "  - NOTE: ~/shell/rc/custom/ and ~/shell/custom-bin/ are inside ~/shell and go with it"
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
# ~/.tmux is legacy - CPS cloned TPM there before tmux dropped its plugins
rm -rf "$HOME/.tmux"
rm -rf "$HOME/.config/nvim"
rm -rf "$HOME/.config/cps"

echo "==> removing Neovim caches and state"
rm -rf "$HOME/.local/share/nvim"
rm -rf "$HOME/.local/state/nvim"
rm -rf "$HOME/.cache/nvim"

echo "==> removing app state outside ~/shell"
rm -rf "$HOME/.local/share/code-server"
rm -rf "$HOME/.config/code-server"
rm -rf "$HOME/.local/kitty.app"
rm -rf "$HOME/.cache/uv"

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
rm -rf "$HOME/.config/lsd"
rm -f "$HOME/.config/starship.toml"

if command -v brew >/dev/null 2>&1; then
  echo "==> uninstalling CPS-installed brew formulas"
  # Kept (installed by cps prereq, broadly useful): wget zip unzip file tmux htop
  # neovim is legacy here - CPS installed it via brew before it became an app bundle.
  brew uninstall \
    neovim \
    nmap openssl ffmpeg imagemagick \
    awscli azure-cli \
    2>/dev/null || true

  echo "==> uninstalling CPS-installed brew casks"
  # Kept (desktop apps a user keeps using): nikitabobko/tap/aerospace sol font-jetbrains-mono-nerd-font
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
echo "  2. cps prereq && cps shell"

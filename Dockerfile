# syntax=docker/dockerfile:1.7

FROM golang:1.25-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w -X 'github.com/tanq16/cli-productivity-suite/cmd.AppVersion=${VERSION}'" \
      -o /out/cps .

FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    TZ=UTC \
    HOMEBREW_NO_AUTO_UPDATE=1 \
    HOMEBREW_NO_ANALYTICS=1 \
    HOMEBREW_NO_INSTALL_CLEANUP=1 \
    HOMEBREW_NO_ENV_HINTS=1

RUN apt-get update && apt-get install -y --no-install-recommends \
      git curl ca-certificates gnupg \
      zsh tmux sudo less \
      build-essential file procps xz-utils unzip \
      locales tzdata \
    && locale-gen en_US.UTF-8 \
    && update-locale LANG=en_US.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

ARG USERNAME=cps
ARG USER_UID=1000
ARG USER_GID=1000
# Ubuntu 24.10+ base ships a default 'ubuntu' user at UID 1000; remove it so cps can claim that UID.
RUN userdel -r ubuntu 2>/dev/null || true \
    && groupdel ubuntu 2>/dev/null || true \
    && groupadd -g ${USER_GID} ${USERNAME} \
    && useradd -m -u ${USER_UID} -g ${USER_GID} -s /usr/bin/zsh ${USERNAME} \
    && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

COPY --from=builder --chown=${USERNAME}:${USERNAME} /out/cps /home/${USERNAME}/.local/bin/cps
RUN chmod +x /home/${USERNAME}/.local/bin/cps

USER ${USERNAME}
WORKDIR /home/${USERNAME}
ENV PATH=/home/${USERNAME}/.local/bin:$PATH

RUN NONINTERACTIVE=1 /bin/bash -c \
      "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

ENV PATH=/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:/home/${USERNAME}/.local/bin:$PATH

# Per-RUN cleanup keeps caches from being baked into the layer (Docker layers can't truly delete).
RUN cps init \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache

# Bake .zprofile so subsequent `zsh -l -c` RUN steps source cps rc fragments.
# Early-returns for interactive logins so `docker exec -it ... zsh -l` doesn't double-source
# (.zshrc handles interactive shells). Not deployed on hosts — purely a Docker build helper.
RUN cat > ~/.zprofile <<'EOF'
[[ -o interactive ]] && return
for f in "$HOME/shell/rc/"*.zsh(N); do source "$f"; done
for f in "$HOME/shell/rc/custom/"*.zsh(N); do source "$f"; done
EOF

SHELL ["/usr/bin/zsh", "-l", "-c"]

RUN cps extend essentials \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend core \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend runtimes \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend cloud \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend security \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend cloudsec \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend appsec \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend misc \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20

RUN cps download-known-extensions && sleep 20
RUN cps extend ai-tools \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend additional-cloud-tools \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend database \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20
RUN cps extend praetorian \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache \
 && sleep 20

# Truly-private tools (toon, nblm, cybernest, lincli) need --gh-token; skipped here.
RUN cps extend private nits raikiri gcli box claudex \
 && sudo rm -rf ~/.cache ~/shell/npm-cache ~/shell/go/cache

CMD ["sleep", "infinity"]

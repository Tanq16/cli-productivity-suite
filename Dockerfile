# syntax=docker/dockerfile:1.7

# =============================================================================
# Stage 1 — build the cps binary
# =============================================================================
FROM golang:1.25-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w -X 'github.com/tanq16/cli-productivity-suite/cmd.AppVersion=${VERSION}'" \
      -o /out/cps .

# =============================================================================
# Stage 2 — Ubuntu sandbox with full cps environment baked in
# =============================================================================
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

# Base apt packages: prerequisites for brew + cps init + day-to-day shell work.
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
# Ubuntu 24.10+ base images ship a default 'ubuntu' user at UID/GID 1000;
# remove it so we can claim that UID for the cps user (matches host UID on
# most Linux dev machines, so bind-mounted volumes Just Work).
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

# Linuxbrew. NONINTERACTIVE skips the "press enter to continue" prompt; sudo
# NOPASSWD lets the install script grab /home/linuxbrew without a TTY.
RUN NONINTERACTIVE=1 /bin/bash -c \
      "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" \
    && /home/linuxbrew/.linuxbrew/bin/brew --version >/dev/null \
    && /home/linuxbrew/.linuxbrew/bin/brew shellenv >/dev/null

# Brew on PATH for the remaining RUN steps (cps init/extend invokes brew).
ENV PATH=/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:/home/${USERNAME}/.local/bin:$PATH

# Phase 1 — base shell env, neovim/nvchad, tmux/TPM, plugin clones, rc fragments.
RUN cps init

# Use login zsh for all subsequent RUN steps so ~/.zshrc sources the cps rc
# fragments — env vars (UV_TOOL_DIR, BUN_INSTALL, FNM_DIR, ...) get set the
# same way an interactive shell would, and cps subprocesses inherit them.
SHELL ["/usr/bin/zsh", "-l", "-c"]

RUN cps extend essentials && sleep 20
RUN cps extend core && sleep 20
RUN cps extend runtimes && sleep 20
RUN cps extend cloud && sleep 20
RUN cps extend security && sleep 20
RUN cps extend cloudsec && sleep 20
RUN cps extend appsec && sleep 20
RUN cps extend misc && sleep 20

# Phase 3 — reference custom-extension packs from the CPS repo.
RUN cps download-known-extensions && sleep 20
RUN cps extend ai-tools && sleep 20
RUN cps extend additional-cloud-tools && sleep 20
RUN cps extend database && sleep 20
RUN cps extend praetorian && sleep 20

# Phase 4 — public-repo tools from the 'private' pack. The truly-private ones
# (toon, nblm, cybernest, lincli — marked IsPrivate: true in the manifest) are
# skipped since they need --gh-token.
RUN cps extend private nits raikiri gcli box claudex

# Sandbox stays alive; user exec's in:  docker exec -it <name> zsh -l
CMD ["sleep", "infinity"]

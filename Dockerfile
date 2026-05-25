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

RUN cps init

# Login zsh so ~/.zprofile sources cps rc fragments → env vars available to cps subprocesses.
SHELL ["/usr/bin/zsh", "-l", "-c"]

RUN cps extend essentials && sleep 20
RUN cps extend core && sleep 20
RUN cps extend runtimes && sleep 20
RUN cps extend cloud && sleep 20
RUN cps extend security && sleep 20
RUN cps extend cloudsec && sleep 20
RUN cps extend appsec && sleep 20
RUN cps extend misc && sleep 20

RUN cps download-known-extensions && sleep 20
RUN cps extend ai-tools && sleep 20
RUN cps extend additional-cloud-tools && sleep 20
RUN cps extend database && sleep 20
RUN cps extend praetorian && sleep 20

# Truly-private tools (toon, nblm, cybernest, lincli) need --gh-token; skipped here.
RUN cps extend private nits raikiri gcli box claudex

RUN sudo rm -rf \
      "$HOME/.cache/Homebrew" \
      "$HOME/.cache/uv" \
      "$HOME/.cache/npm" \
      "$HOME/.cache/go-build" \
      "$HOME/shell/npm-cache" \
      "$HOME/shell/go/cache"

CMD ["sleep", "infinity"]

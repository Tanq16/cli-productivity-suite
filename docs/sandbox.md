# The sandbox container

`tanq16/cps-sandbox` is an Ubuntu image with a fully built CPS environment inside it. It exists for machines where installing CPS directly is not wanted: a throwaway workspace, a CI box, or a container you hand to an AI coding agent.

```bash
docker run -d --name cps-sandbox tanq16/cps-sandbox:latest
docker exec -it cps-sandbox zsh -l
```

The default command is `sleep infinity`, so the container stays up and you `docker exec` in whenever you need it. The `zsh -l` is what sources the rc fragments, so it is what gives you the configured shell with PATH wired, plugins loaded and the starship prompt. Inside, `tt` starts a tmux session and `t` re-attaches. Add `--rm` to `docker run` for an ephemeral run, or keep the container and re-`exec` for one you come back to.

## What is inside

Every extension pack is pre-installed, so nothing bootstraps at first use. Only the two token-gated tools in the `private` pack, `toon` and `cybernest`, are missing, since they need `--gh-token` at build time.

The user is a non-root `cps` at UID and GID 1000, with passwordless sudo for ad-hoc installs. Homebrew is installed as Linuxbrew.

That combination is what makes the image useful for an agent: `claude`, `codex`, `cursor-agent` and `agy` are already on PATH, alongside every runtime, search tool, cloud CLI and scanner they are likely to reach for, and all of it runs as a non-root user inside a disposable container rather than on your host.

```bash
docker run -d --name agent-sandbox tanq16/cps-sandbox:latest
docker exec -it agent-sandbox zsh -l
# inside:
claude
```

The image is multi-arch (`linux/amd64` and `linux/arm64`) and several GB.

## Building it

```bash
docker build -t cps-sandbox .
docker run -d --name cps-sandbox cps-sandbox
docker exec -it cps-sandbox zsh -l
```

Published builds come from the Docker Sandbox workflow, which is manually dispatched and takes an image tag defaulting to `latest`. It pushes both that tag and the commit SHA.

## Verifying a build

```bash
docker run --rm -v "$PWD/scripts/verify.sh:/tmp/v.sh" \
    tanq16/cps-sandbox:latest zsh -lc 'bash /tmp/v.sh'
```

`verify.sh` asserts the user and sudo setup, the `~/shell` tree, the rc fragments, environment variables, PATH order, and one binary per pack. It stays silent on success and exits 1 listing whatever is missing. The `zsh -l` wrapper is required, because that is what sources the rc fragments it checks.

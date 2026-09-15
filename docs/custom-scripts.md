# Custom scripts

`cps custom` runs `~/.config/cps/custom.sh` and gets out of the way. CPS resolves the path, execs the script with the CPS `PATH` in front, passes its arguments through verbatim, and exits with whatever the script exited with. Unknown names, argument parsing and error messages all belong to the script.

Flag-shaped arguments go after a `--` separator, so that `cps custom --debug` stays the CPS debug flag and `cps custom -- --debug` reaches the script. Bare names need no separator:

```bash
cps custom list
cps custom mytool
cps custom -- --help
```

The script must be executable:

```bash
chmod +x ~/.config/cps/custom.sh
```

It is run directly rather than through a shell, so its shebang decides the interpreter.

## The environment CPS provides

`PATH` is prefixed with these, in order:

```
~/shell/custom-bin
~/shell/extensions
~/shell/go-sdk/bin
~/shell/go/bin
~/shell/fnm/aliases/lts-latest/bin
~/shell/uv-tool-executables
```

Everything else is inherited from the calling shell.

## The contract

1. `custom.sh list` prints one package name per line.
2. `custom.sh <name>` installs or updates that package.
3. Exit 0 means success. Any non-zero exit means failure, and `cps custom` returns the same code.
4. The script never writes `~/.config/cps/state.json`. It reads installed state through `cps status --json`.
5. Shell fragments go in `~/shell/rc/custom/*.zsh` and binaries in `~/shell/custom-bin/`.

Rule 4 is the one that matters. Two writers of `state.json` with no schema guard is what made the retired YAML pack system fragile, and `cps status --json` is the supported way to read that state:

```bash
cps status --json | jq -r '.[] | select(.installed) | .name'
```

## Example

```bash
#!/usr/bin/env bash
set -euo pipefail

PACKAGES=(mytool otherthing)

case "${1:-}" in
  list)
    printf '%s\n' "${PACKAGES[@]}"
    ;;
  mytool)
    curl -fsSL https://example.com/mytool -o ~/shell/custom-bin/mytool
    chmod +x ~/shell/custom-bin/mytool
    ;;
  otherthing)
    echo 'export OTHER_HOME="$HOME/otherthing"' > ~/shell/rc/custom/otherthing.zsh
    ;;
  *)
    echo "usage: cps custom {list|${PACKAGES[*]}}" >&2
    exit 1
    ;;
esac
```

## What CPS will not do

CPS does not track, version, update or remove anything the script installs. `cps status` reports CPS packages only, and `./scripts/deep-removal.sh` leaves `~/shell/custom-bin/` and `~/shell/rc/custom/` alone.

# App bundles

Three tools in the `homelab` pack are multi-file trees rather than single binaries, so they unpack to `~/shell/apps/<name>/` and are not on PATH. Run them by full path.

```bash
~/shell/apps/rinnegan/bin/rinnegan serve            # PTY web terminal
~/shell/apps/code-server/bin/code-server            # VS Code in the browser, http://127.0.0.1:8080
~/shell/apps/neo4j/bin/neo4j console                # graph database, http://localhost:7474
~/shell/apps/neo4j/bin/cypher-shell                 # Cypher client for the above
```

An upgrade replaces the whole tree, so nothing of yours may live inside a bundle. Each of these keeps its state elsewhere for exactly that reason.

## code-server

CPS seeds the config on first install and never rewrites it, so anything you change afterwards sticks. `~/.config/code-server/config.yaml` binds it to `127.0.0.1:8080` with authentication disabled, telemetry and update checks off, workspace trust off, and the port-proxy routes disabled. Editor settings live in `~/.local/share/code-server/User/settings.json`, themed Catppuccin Mocha with JetBrains Mono.

Every key in that YAML is a code-server flag with the leading dashes stripped, and an unknown key is a fatal startup error. CLI flags beat the config file, so a different port for one launch is `code-server --bind-addr 127.0.0.1:9000`.

`auth: none` is only safe because it binds to loopback. To reach it from another machine, forward the port over SSH rather than changing `bind-addr`:

```bash
ssh -L 8080:localhost:8080 you@host
```

That also keeps the browser on `localhost`, which is a secure context. Over a plain-HTTP LAN address, webviews and clipboard access break.

To add authentication, replace `auth: none` with `auth: password` plus either `password:` or `hashed-password:`, the argon2 hash winning if both are set. Neither can be passed as a CLI flag by design, so they go in the config file or in the `PASSWORD` and `HASHED_PASSWORD` environment variables.

Eight extensions are installed on first seeding, from Open VSX: Catppuccin theme and icons, EditorConfig, Error Lens, Go, Python, Ruff, and Claude Code. Seeding only happens while the extensions directory is empty, so adding or removing one by hand sticks, and a failed install retries on the next `cps extend homelab`.

The Go extension will offer to install `gopls` on first use, which needs `cps extend runtimes`. Pylance is not on Open VSX, so Python IntelliSense is weaker than desktop VS Code and Ruff covers linting and formatting. The Claude Code panel runs on top of the `claude` CLI and shares its login, so it uses an existing Claude subscription rather than a metered API key. Install that CLI with `cps extend ai-tools` and authenticate once on the host.

## neo4j

The bundle ships no JVM. It runs on the `JAVA_HOME` that `cps extend runtimes` sets up, which is why `homelab` is marked as needing that pack.

Neo4j would otherwise keep its databases inside the bundle, where an upgrade would destroy them. `cps extend homelab` seeds `~/.config/neo4j/conf/` once, never overwriting it afterwards, with absolute `server.directories.*` paths pointing at `~/.config/neo4j/`, and `50-homelab.zsh` exports `NEO4J_CONF` so Neo4j reads that conf instead of the bundle's. Your graph, plugins and tuning survive every upgrade, and `scripts/deep-removal.sh` leaves them alone.

Set a password before the first start:

```bash
~/shell/apps/neo4j/bin/neo4j-admin dbms set-initial-password <password>
```

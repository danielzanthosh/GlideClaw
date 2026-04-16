# GlideClaw

GlideClaw is a lightweight, terminal-first personal gateway agent for Linux servers.

It is designed for low-resource hosts (small VPS, Raspberry Pi class machines) and focuses on:
- local/headless operation,
- strict command policy boundaries,
- Telegram as a primary remote surface,
- modular connectors,
- SQLite-backed state.

## Project status (honest)

GlideClaw is in **early runtime scaffold** stage.

### Implemented now
- Single binary Go CLI (`glideclaw`).
- Config loading with defaults, path resolution, and env overrides.
- `init`, `status`, `run`, `doctor`, and security/password CLI scaffolding.
- SQLite schema bootstrap and local state initialization.
- Tiered policy engine + Tier 3 escalation/password gates.
- Connector registry with **placeholder/not-configured** connector stubs.
- Systemd service example and architecture docs.

### Planned / incomplete
- Full Telegram message handling lifecycle.
- Real connector auth flows and API operations (Google, GitHub, Vercel).
- Full memory retrieval/summarization behavior.
- Archive lifecycle automation and restore UX.

## Architecture summary

- `cmd/glideclaw`: entrypoint.
- `internal/config`: path resolution + config loading.
- `internal/setup`: bootstrap/init workflow.
- `internal/db`: SQLite schema/migration bootstrap.
- `internal/policy` + `internal/executor`: command classification + execution gates.
- `internal/security`: escalation password and lockout state.
- `internal/connectors`: registry and connector scaffolds.
- `internal/telegram`: adapter scaffold.

Detailed blueprint: `docs/ARCHITECTURE.md`.

## Build

```bash
go build -o bin/glideclaw ./cmd/glideclaw
```

## Run and onboarding

> **Security warning:** Do not run GlideClaw as root for normal deployments. Use a dedicated service account with restricted filesystem access.

1. Initialize local state:

```bash
./bin/glideclaw init
```

2. Review generated config/bootstrap files.
3. Set escalation password:

```bash
./bin/glideclaw security set-password
```

4. Start daemon scaffold:

```bash
./bin/glideclaw run
```

Check state:

```bash
./bin/glideclaw status
```

## Commands currently available

- `glideclaw init`
- `glideclaw status`
- `glideclaw run`
- `glideclaw doctor`
- `glideclaw config validate`
- `glideclaw connector status`
- `glideclaw archive run`
- `glideclaw security set-password`
- `glideclaw security change-password`
- `glideclaw security status`
- `glideclaw security reset-lockout`
- `glideclaw exec [--override-safe] <command>`

## Default path behavior

- **Non-root mode:**
  - config: `~/.config/glideclaw/config.yaml`
  - bootstrap: `~/.config/glideclaw/BOOTSTRAP.md`
  - data: `~/.local/share/glideclaw/`
- **Root/system mode (not recommended for dev):**
  - config: `/etc/glideclaw/config.yaml`
  - bootstrap: `/etc/glideclaw/BOOTSTRAP.md`
  - data: `/var/lib/glideclaw/`

## Repository layout

```text
cmd/               # binary entrypoint
internal/          # core runtime packages
configs/           # config examples
docs/              # architecture + deployment docs
systemd/           # service unit examples
```

## Developer quick start

```bash
go mod tidy
go test ./...
go run ./cmd/glideclaw init
go run ./cmd/glideclaw status
```

## Connectors note

Google, GitHub, and Vercel connectors are currently scaffold placeholders and are not fully implemented integrations yet.

## License

Apache License 2.0. See `LICENSE` and `NOTICE`.

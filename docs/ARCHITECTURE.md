# GlideClaw Master Blueprint (Implementation Grade)

## 1) Product architecture

### Module map
- `cmd/glideclaw`: single binary entrypoint.
- `internal/app`: dependency wiring and runtime boot.
- `internal/cli`: terminal-first command router and interactive loop.
- `internal/telegram`: Telegram adapter (pairing, DM/group routing, file intake).
- `internal/db`: SQLite bootstrap and migrations.
- `internal/bootstrap`: `BOOTSTRAP.md` parser and policy hints.
- `internal/policy`: command risk engine + safe mode gates.
- `internal/archive`: hot/cold lifecycle + Drive offload orchestration.
- `internal/connectors`: connector interfaces, registry, auth boundaries.

### Core runtime diagram (text)
1. `glideclaw run` starts gateway process.
2. Config + BOOTSTRAP are loaded; SQLite opened (WAL mode).
3. Connector registry initializes lazily (no eager heavy clients).
4. Telegram loop receives update or CLI receives command.
5. Message enters agent pipeline: authorize -> memory fetch -> policy evaluation -> tool/connector action.
6. Action and result are persisted (messages, audit, approvals, archive metadata).
7. Background maintenance ticks at low frequency: retention sweep, archive candidate scan, token health check.

### Responsibilities
- **Gateway core**: surface-agnostic message/event handling.
- **Policy engine**: strict risk tiers, safe mode override, workspace restrictions.
- **Memory service**: working/session/long-term/archive metadata separation.
- **Archive manager**: local-first cache, offload and restore with checksum verification.
- **Connectors**: isolated auth and failure domains.

## 2) Folder structure

```text
.
├── cmd/glideclaw/main.go
├── configs/config.example.yaml
├── docs/ARCHITECTURE.md
├── BOOTSTRAP.example.md
├── systemd/glideclaw.service
└── internal
    ├── app
    ├── archive
    ├── bootstrap
    ├── cli
    ├── config
    ├── connectors
    │   ├── google
    │   ├── github
    │   └── vercel
    ├── db
    ├── policy
    └── telegram
```

- Keep each subsystem in its own package to keep compile units small and testable.
- Core binary remains single-process, no sidecar requirement.

## 3) Config schema

See `configs/config.example.yaml` for a complete template. Features:
- runtime profile (`micro|lite|normal`)
- gateway + Telegram + terminal toggles
- connector enablement and scope lists
- archive thresholds and retention
- execution and security policy sections
- logging and task scheduler knobs

Env overrides are prefixed with `GLIDECLAW_` (e.g., `GLIDECLAW_TELEGRAM_BOT_TOKEN`).

## 4) BOOTSTRAP.md spec

Markdown-first contract:
- section heading (`## Section`) + bullet values.
- unknown sections preserved in `RawSections` for forward compatibility.
- parsed hints directly influence response style, risk boundaries, and memory weighting.

Supported sections:
- Identity
- Preferences
- Environments
- Project defaults
- Allowed autonomous actions
- Confirmation-required actions
- Blocked actions
- Security mode
- Connector notes
- Deployment preferences
- Memory hints

## 5) SQLite schema

Implemented in `internal/db/schema.go`.

Tables:
- `sessions`, `messages`
- `memory_entries`, `memory_pins`
- `tasks`
- `connector_auth`
- `archive_objects`
- `file_metadata`
- `command_approvals`
- `audit_log`

Design points:
- WAL mode + single connection for low-RAM stability.
- index coverage for retrieval-heavy paths.
- metadata-only archive pointers keep local state tiny.

## 6) Connector system design

Connector contract (`internal/connectors/interface.go`):
- `Name()`
- `Enable(ctx)` / `Disable(ctx)`
- `Health(ctx)`

Rules:
- independently enabled and audited.
- separate auth record per connector (`connector_auth`).
- lazy instantiate API clients only on first operation.
- connector errors return bounded failure, never process crash.

Auth model:
- Google: OAuth device flow per product scope set (Drive/Gmail/Calendar separate).
- GitHub: fine-grained PAT now; app auth later.
- Vercel: scoped team/project token.

## 7) Google Drive archive subsystem

### Hot/cold flow
1. File enters hot cache (`archive.hot_dir`).
2. Metadata + checksum written locally.
3. Sweep marks cold candidates by age/size profile.
4. Upload to Drive archive folder, persist `drive_file_id`.
5. Local large blob evicted; metadata retained.

### Restore on demand
- request object id -> fetch from Drive -> verify SHA-256 -> place in restore cache -> set TTL.
- corrupted checksum triggers quarantine state and audit record.

### Eviction strategy
- profile-dependent TTL and size cap.
- LRU + oldest-first for restore cache.
- dry-run mode prints decisions only.

## 8) Command security system

Policy tiers:
- **Tier 0** always allowed readonly.
- **Tier 1** policy-allowed bounded writes/build/test.
- **Tier 2** explicit approval required.
- **Tier 3** blocked by default.

Engine controls:
- command deny prefixes + BOOTSTRAP blocked phrases.
- workspace allowlist and timeout caps.
- safe mode: non-tier0 downgraded to require explicit approval.
- every request records audit + optional `command_approvals` entry.

## 9) Telegram design

Lifecycle:
1. Admin sets bot token.
2. Unknown user sends `/pair` -> pending pairing record.
3. Local admin approves via CLI (`pair approve`).
4. DM chat enabled with context continuity.
5. Group mode opt-in with explicit allowlist.

Safety:
- per-chat rate limits.
- max attachment size gate.
- MIME allowlist and archive pipeline integration.
- admin commands segregated (`/safe_mode`, `/connector_status`, `/archive_run`).

## 10) CLI design

Core command families:
- bootstrap/config: `init`, `config validate`, `bootstrap show/edit`
- runtime: `run`, `doctor`, `logs`, `chat`
- auth/pairing: `pair list|approve`
- memory/archive: `memory add/search/pin`, `archive status/run/restore`
- connectors: `connect ...`, `connector status`
- governance: `exec policy`, `safe-mode on/off`
- tasks/backup: `task add/list/run`, `backup create/restore`

UX principles:
- terse defaults, `--json` optional later.
- deterministic exit codes for automation.

## 11) systemd deployment

See `systemd/glideclaw.service`.

Highlights:
- `Type=simple`
- restart with backoff
- locked-down service settings (`NoNewPrivileges`, `PrivateTmp`, strict umask)
- data under `/var/lib/glideclaw`, logs via journald + JSONL file

## 12) Phased roadmap

### MVP
- single daemon + CLI skeleton
- Telegram DM + pairing
- SQLite schema + memory add/search basics
- policy engine with safe mode
- archive metadata and dry-run offload

### v1
- full connector auth flows
- task scheduler and recurring jobs
- archive restore with checksum verification
- richer CLI/TUI chat loop + approvals queue

### Advanced
- semantic memory option (pluggable embeddings)
- service-account Google mode
- GitHub app mode and repo write workflows
- Vercel deploy approvals and policy packs

### Future optional
- multi-agent task routing
- end-to-end encrypted attachment vault
- remote backup target alternatives (S3/B2)

## 13) Production scaffold status

Implemented in this repository:
- main entrypoint
- package skeletons
- config loader + env overrides
- SQLite bootstrap schema
- CLI router skeleton
- Telegram adapter skeleton
- connector interface + registry + stubs
- archive manager skeleton
- policy engine skeleton

This baseline is deliberately small and safe so it can run on weak Linux hosts, then evolve iteratively without architectural rewrites.

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
- `internal/executor`: policy-integrated command runner with escalation checks.
- `internal/security`: password hashing, lockout, and elevation windows.
- `internal/audit`: Tier 3 audit logger.
- `internal/archive`: hot/cold lifecycle + Drive offload orchestration.
- `internal/connectors`: connector interfaces, registry, auth boundaries.

### Core runtime diagram (text)
1. `glideclaw run` starts gateway process.
2. Config + BOOTSTRAP are loaded; SQLite opened (WAL mode).
3. Connector registry initializes lazily (no eager heavy clients).
4. Telegram loop receives update or CLI receives command.
5. Message enters agent pipeline: authorize -> memory fetch -> policy evaluation -> tool/connector action.
6. Tier 3 flow: classify risk -> check hard block -> safe-mode override gate -> password auth (Argon2id/PBKDF2) -> optional typed double confirmation -> execute or deny.
7. All Tier 3 outcomes are persisted to `audit_log` without passwords/hashes.

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
    ├── audit
    ├── bootstrap
    ├── cli
    ├── config
    ├── connectors
    ├── db
    ├── executor
    ├── policy
    ├── security
    └── telegram
```

## 3) Config schema

See `configs/config.example.yaml` for a full template.

Escalation fields:
- `security.escalation_enabled`
- `security.elevation_mode` (`single|time_window`)
- `security.elevation_window_seconds` (default 60)
- `security.max_attempts`
- `security.lockout_seconds`
- `security.require_double_confirmation`
- `security.critical_confirm_text`
- `security.allow_tier3_in_safe_mode`

## 4) BOOTSTRAP.md spec

Markdown-first contract remains unchanged. `Blocked actions` continue to hard-block execution even with valid escalation password.

## 5) SQLite schema

Implemented in `internal/db/schema.go`.

Audit requirement is met via `audit_log` entries for every Tier 3 attempt:
- timestamp (`created_at`)
- command attempted (`target`)
- source + execution result (`details_json`)
- result (`outcome`)

## 6) Connector system design

No privilege change: escalation only authorizes command execution layer. Connector scopes and auth boundaries still apply.

## 7) Google Drive archive subsystem

Unchanged hot/cold architecture; restore/eviction still governed by archive policy.

## 8) Command security + Tier 3 escalation design

### Risk model
- **Tier 0** allowlisted readonly commands.
- **Tier 1** normal policy-allowed commands.
- **Tier 3 escalation** dangerous-but-allowable commands requiring auth.
- **Hard blocked** commands are denied always.

### Escalation credentials
- Password set by owner via `glideclaw security set-password`.
- Password hash storage in local secrets dir with `0600` mode.
- Uses **Argon2id** by default, with **PBKDF2** fallback option.
- Unique random 16-byte salt per password set.
- Never stores plaintext password.

### Lockout + brute-force controls
- Failed attempts counted locally.
- Lockout after `max_attempts` for `lockout_seconds`.
- Lockout state persisted in secrets dir with `0600`.

### Elevation window
- `single`: one command after successful auth.
- `time_window`: valid for `elevation_window_seconds`.

### Safe mode
- If safe mode is on, Tier 3 is denied unless `--override-safe` is supplied and policy allows safe-mode Tier 3 override.

### Double confirmation
- For critical actions, require typed confirmation token (e.g., `DELETE_PRODUCTION_DATA`) after password verification.

## 9) Telegram design (escalation-specific)

- Tier 3 from Telegram requires secure challenge-response flow.
- Password replies are ephemeral and must not be persisted as normal chat messages.
- Approved/denied events still written to audit log.

## 10) CLI design

Added security commands:
- `glideclaw security set-password`
- `glideclaw security change-password`
- `glideclaw security status`
- `glideclaw security reset-lockout`

Execution command:
- `glideclaw exec <command>`
- `glideclaw exec --override-safe <command>`

## 11) systemd deployment

No model changes required. Keep secret files under `security.secrets_dir` with service `UMask=0077`.

## 12) Phased roadmap

### MVP refinement
- complete Tier 3 escalation flow in terminal + audit + lockout.

### v1
- Telegram secure escalation UX and ephemeral password handling.

### Advanced
- optional external secret backends (pass/KMS/HSM wrappers).

## 13) Production scaffold status

Now includes implementation scaffolding for:
- Argon2id/PBKDF2 password hashing
- lockout/elevation manager
- Tier 3 audit logger
- executor-policy integration and CLI security commands

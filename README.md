# GlideClaw

GlideClaw is a lightweight, terminal-first, Telegram-first personal gateway agent for low-end Linux servers.

- single Go binary
- SQLite-backed tiered memory
- policy-gated command execution
- connector architecture (Google Workspace, GitHub, Vercel)
- archive layer using local cache + Google Drive cold storage

## Quick start

```bash
go run ./cmd/glideclaw doctor
```

Read the full implementation blueprint in `docs/ARCHITECTURE.md`.

## License

This project is licensed under the Apache License 2.0.
See the LICENSE file for details.

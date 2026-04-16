# Deployment Guide (Scaffold Stage)

This document describes a practical non-root deployment for the current GlideClaw scaffold.

## 1. Build

```bash
go build -o glideclaw ./cmd/glideclaw
sudo install -m 0755 glideclaw /usr/local/bin/glideclaw
```

## 2. Create service account and directories

```bash
sudo useradd --system --create-home --home-dir /var/lib/glideclaw --shell /usr/sbin/nologin glideclaw
sudo mkdir -p /etc/glideclaw /var/lib/glideclaw
sudo chown -R glideclaw:glideclaw /var/lib/glideclaw
sudo chmod 0750 /var/lib/glideclaw
```

## 3. Initialize runtime state

Run as the glideclaw user:

```bash
sudo -u glideclaw /usr/local/bin/glideclaw init
```

Then review:
- config file,
- BOOTSTRAP profile,
- secret placeholders.

## 4. Configure environment

Create `/etc/glideclaw/glideclaw.env` with minimal required values (example placeholders only):

```bash
GLIDECLAW_PROFILE=lite
GLIDECLAW_SAFE_MODE=true
# GLIDECLAW_TELEGRAM_BOT_TOKEN=<set-if-enabled>
```

Restrict permissions:

```bash
sudo chown root:glideclaw /etc/glideclaw/glideclaw.env
sudo chmod 0640 /etc/glideclaw/glideclaw.env
```

## 5. Systemd

Copy unit:

```bash
sudo cp systemd/glideclaw.service /etc/systemd/system/glideclaw.service
sudo systemctl daemon-reload
sudo systemctl enable --now glideclaw
```

## 6. Inspect status and logs

```bash
sudo systemctl status glideclaw
journalctl -u glideclaw -f
sudo -u glideclaw /usr/local/bin/glideclaw status
```

## 7. Non-root guidance

- Prefer dedicated non-root user.
- Avoid broad write permissions outside GlideClaw data/config dirs.
- Keep `safe_mode` enabled unless you have a controlled reason.

## 8. Backup expectations

At minimum back up:
- SQLite DB file,
- BOOTSTRAP profile,
- config file,
- secrets directory metadata.

Do **not** back up plaintext tokens into public or shared locations.

## 9. Connector caution

Connectors are currently scaffolded placeholders. Treat related auth wiring as incomplete until connector implementation hardening is complete.

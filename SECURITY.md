# Security Policy

GlideClaw executes local commands and stores connector-related authentication state. Treat deployments as security-sensitive.

## Reporting vulnerabilities

Please report vulnerabilities privately before public disclosure.

- Preferred: open a **private security report** (or contact maintainers directly if a private channel is listed by the project).
- Include: affected version/commit, impact, reproduction steps, and suggested remediation if available.

Do **not** open public issues for unpatched critical vulnerabilities.

## Supported versions

At this stage, only the latest commit on the default branch is considered supported for security fixes.

## Threat warning

GlideClaw can:
- execute shell commands,
- manage automation actions,
- store escalation/auth material,
- access connector credentials.

A misconfigured deployment can create high-impact risk.

## Deployment safety requirements

1. Run under a dedicated **non-root** service account.
2. Restrict filesystem access to explicit config/data directories only.
3. Keep network exposure minimal; do not expose unsafe admin surfaces publicly.
4. Use least-privilege connector scopes and rotate tokens regularly.
5. Keep secrets out of Git and out of plaintext shared channels.

## Secret handling expectations

- Do not commit `.env`, token files, database snapshots, or local secret stores.
- Restrict secret files with owner-only permissions where possible.
- Prefer environment files readable only by the GlideClaw service account.

## Connector scope minimization

Always request the smallest feasible scopes:
- Google: separate Drive/Gmail/Calendar scopes only as needed.
- GitHub: fine-grained permissions over broad tokens.
- Vercel: project/team scoped tokens.

## Disclosure expectations

After receiving a report, maintainers should:
1. Acknowledge receipt quickly.
2. Reproduce and triage severity.
3. Prepare and release a fix.
4. Publish a short advisory/changelog note once users can patch.

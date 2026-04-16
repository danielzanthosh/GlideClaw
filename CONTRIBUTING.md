# Contributing to GlideClaw

Thanks for contributing.

## Local setup

```bash
git clone <repo-url>
cd GlideClaw
go mod tidy
```

## Build

```bash
go build -o bin/glideclaw ./cmd/glideclaw
```

## Test

```bash
go test ./...
```

## Development workflow

1. Create a feature branch.
2. Keep changes focused and modular.
3. Run formatting/tests before opening PR.
4. Document behavior changes in README/docs where relevant.

## Coding style

- Keep dependencies minimal.
- Prefer small packages and clear interfaces.
- Avoid overengineering; preserve low-resource design goals.
- Keep comments practical and concise.

## Security expectations

- Never commit real secrets, tokens, or local runtime state.
- Treat command execution and connector auth as sensitive areas.
- Follow least privilege for any new connector scopes.
- If you discover a vulnerability, follow `SECURITY.md`.

## Pull requests

PRs should include:
- clear summary and motivation,
- testing commands + results,
- any migration or operational notes,
- accurate claims (mark scaffolded/placeholder features clearly).

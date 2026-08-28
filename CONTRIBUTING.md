# Contributing to Project PULSE

Thanks for your interest in Project PULSE!
This guide covers what you need to know before opening a PR.

## Before you start: the CLA

This project is distributed under a **Business Source License 1.1, non-commercial** (see [`ADR-005`](docs/adr/ADR-005-source-available-licence.md)). Any external contribution (code, docs, etc.) requires accepting the [**CLA.md**](CLA.md) beforehand. Without a signed CLA, a PR cannot be merged, no matter how good it is technically.

## Development environment

The simplest path is the **devcontainer** (`.devcontainer/`) via GitHub Codespaces or VS Code ("Reopen in Container") — Go, Node, `sqlc`, `air`, `golang-migrate`, and `postgresql-client` are provisioned automatically, and `docker compose up -d postgres redis` starts the dependencies.

For local setup without the devcontainer, see the prerequisites in [`README.md`](README.md).

## Before every commit: `pre-commit`

```bash
pip install pre-commit
pre-commit install
```

This installs the hooks that run automatically on every `git commit`:

| Hook                    | Purpose                                                                                                                     |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `golangci-lint`         | Go linting (staticcheck, revive, errcheck, govet, ineffassign, unused, misspell)                                            |
| `go-fmt` / `go-imports` | Go formatting                                                                                                               |
| `go-unit-tests`         | Go unit tests                                                                                                               |
| `eslint`                | TypeScript/React linting                                                                                                    |
| `prettier`              | JS/TS/CSS/JSON formatting (not Markdown — see `markdownlint`)                                                               |
| `markdownlint`          | Markdown style (bullets use `*`, see `.markdownlint.json`)                                                                  |
| `sqlfluff-lint`         | SQL linting for migrations (PostgreSQL dialect)                                                                             |
| ADR-003 guard           | Blocks any `JOIN core.\|tournament.\|scheduling.\|finance.\|evaluation.` in migrations — schema isolation is non-negotiable |

A hook that **modifies** files (Prettier, `go fmt`) intentionally fails the commit the first time — review the diff, `git add`, commit again.

To run everything against the whole repo without committing:

```bash
pre-commit run --all-files
```

## Tests

A PR will not be merged if CI (lint, tests, Snyk scan, SonarQube Cloud) is red.

```bash
# Backend
go test -v -race -coverprofile=coverage-go.out ./internal/... ./pkg/...

# Frontend
cd frontend && npm run test -- --coverage
```

Any new HTTP route or new repository must ship with tests. Mocks of `ports.UserRepository` (see `user_handler_test.go`) must honor the same error contracts as the real implementation (e.g. `ports.ErrUserNotFound`), not just satisfy the interface signature.

## Branch naming convention

```text
feature/<short-description>
fix/<short-description>
```

## Architecture Decision Records (ADRs)

Any structural change (a major new dependency, a stack change, a module split, etc.) must come with an ADR in `docs/adr/`, following the pattern of the existing ADRs (context → decision → consequences). A PR that changes a decision already recorded in an ADR must first propose a replacement ADR.

## Code style

* **Go:** must comply with `.golangci.yml` — don't disable a rule without a written justification in the config file.
* **TypeScript/React:** must comply with `frontend/eslint.config.js`. An identifier prefixed with `_` (e.g. `_size`) marks intentional scaffolding not yet wired up — that's not an error, it documents an intent still to be finished.
* **Commits:** explain the *why*, not just the *what*, especially for changes touching security or schema isolation (`core`, `tournament`, `scheduling`, `finance`, `evaluation`).

## Security

Never open a public PR or issue for a vulnerability. See [`SECURITY.md`](SECURITY.md).

# Security Policy

## Scope

This policy covers the Go backend (`cmd/backend`, `internal/`, `pkg/`), the React frontend (`frontend/`), the Dockerfiles (`deployments/docker/`), and the infrastructure configuration (`docker-compose.yml`, `.devcontainer/`) of Project PULSE.

## Supported versions

The project is in active phase-1 development, with no stable release published yet. Security fixes are applied to the default branch only.

| Version                         | Supported |
| ------------------------------- | --------- |
| `main` (latest commit)          | ✅        |
| Any unmaintained fork or branch | ❌        |

This table will be updated once the first tagged release ships.

## Reporting a vulnerability

**Never open a public issue or PR for a vulnerability.**

Recommended channel: use this repository's **Security → Report a vulnerability** tab (GitHub Security Advisories), which enables a private exchange with maintainers without publicly exposing details before a fix is available.

Please include as much of the following as possible:

* a description of the issue and its potential impact
* steps to reproduce it (or a minimal PoC)
* the affected version/commit
* a suggested fix, if you have one

## Areas of particular interest

* Any leak or bypass of PostgreSQL schema isolation (`core`, `tournament`, `scheduling`, `finance`, `evaluation` — see `ADR-003`)
* Authentication/authorization: password hashing, session handling, access control on `/api/v1/*` routes
* Client IP resolution (`middleware.ClientIPFromHeader`) and any spoofing risk relative to the actual deployment network topology
* Leakage of internal details in HTTP error responses
* Vulnerable dependencies (Go and npm) — also covered by the automated Snyk scan in CI

## Out of scope

* Vulnerabilities in third-party dependencies that are already publicly known and for which an upstream fix isn't available yet (please still flag it if you think we should accelerate an update)
* Attacks requiring physical access to a club's self-hosted deployment infrastructure
* Social engineering targeting maintainers or contributors

## What to expect

* Acknowledgment within a reasonable timeframe
* Fixes prioritized based on real-world severity (exploitability, impact on club/player data)
* Credit to the reporter in the fix's changelog, unless they request otherwise

## Security tooling already in place

* **Snyk** dependency scanning in CI (`security-snyk`)
* **SonarQube Cloud** static analysis
* `golangci-lint` with `staticcheck` (catches risky deprecated APIs, e.g. the `middleware.RealIP` CVE)
* **Hadolint** scanning of Dockerfiles
* Production images built on `distroless` (minimal attack surface, no shell)

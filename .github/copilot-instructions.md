# PULSE — AI Agent Instructions

Universal SaaS platform for managing team sports (multi-sport, multi-level, ephemeral rosters, tournaments, dynamic scheduling). Modular monolith in Go designed for a future extraction into microservices.

## Architecture

See `docs/ARCHITECTURE.md`, `docs/ArchitectureTechnique.md`, `docs/ArchitectureC4.md`, and the ADRs in `docs/adr/` (source of truth for any structural decision). Always validate a change against these documents before making it.

* **Pattern**: Hexagonal (Ports & Adapters) per module, inside a modular monolith (`cmd/monolith/main.go`).
* **Modules** (`internal/{module}/`): `core`, `tournament`, `scheduling`, `finance`, `evaluation`. Each follows the structure:

  ```text
  internal/{module}/
  ├── domain/    # pure business logic, no external dependencies
  ├── ports/     # interfaces (input/output contracts)
  └── adapters/postgres/queries/{entity}.sql  # sqlc-style queries
  ```

* **Strict module isolation**: all inter-module communication goes through a Go interface (port). **No SQL JOIN across PostgreSQL schemas is allowed** — one Postgres schema per module (`core`, `tournament`, `scheduling`, `finance`, `evaluation`).
* **Adapter duality**: each module's port must be implementable by an `in-memory` adapter (monolith mode, current default) and, later, a `gRPC/event` adapter (microservice mode) — selected at startup via configuration/env vars, per ADR-003. Never bypass the port to call another module's adapter directly.
* **Multi-sport**: all relevant entities carry a `sport_id` (UUID) to prepare for multi-sport support (Soccer in phase 1).

## Tech stack

* **Backend**: Go 1.25, `go-chi/chi/v5` router, `jackc/pgx/v5` driver, structured logging via `slog` (JSON handler), OpenTelemetry tracing (`riandyrn/otelchi`) and logging (`otelslog`) bridges.
* **Database**: PostgreSQL 16+ (schema per module), Redis 7+ (cache/pub-sub/locks).
* **Frontend**: React + TypeScript, Tailwind CSS, Radix UI/Shadcn — `frontend/` folder still empty (scaffold pending).
* **Containers**: multi-stage Docker builds, non-root distroless images for the backend.

## Conventions

* **SQL migrations**: `migrations/{module}/NNNNNN_description.up.sql` (`golang-migrate` convention), one migration per schema/module.
* **SQL queries**: files under `internal/{module}/adapters/postgres/queries/{entity}.sql`, sqlc-style comments (`-- name: FunctionName :one|:many|:exec`).
* **DB schema**: UUID primary keys (`gen_random_uuid()`), `TIMESTAMP WITH TIME ZONE`, documented via `COMMENT ON TABLE|COLUMN`.
* **Language**: all content (code, identifiers, comments, commit messages, PRs) must be written in English, except the `docs/` directory which stays in French.
* **Accessibility (frontend)**: WCAG 2.1 AA, mobile-first, min. 44x44px touch targets, native Dark/Light support via CSS variables.
* **Licensing**: Business Source License 1.1 (Non-Commercial), see `LICENSE` and ADR-005. Never suggest or add code that would offer PULSE as a paid hosted SaaS to third parties. External contributions require signing `CLA.md`.
* After completing any task, remind the user to update `docs/ROADMAP_TODO.md`.

## Build & run

```bash
cp .env.example .env
docker compose up -d          # backend + PostgreSQL + Redis (+ frontend if present)
curl http://localhost:8080/healthz
```

There is no Makefile, no test suite, and no `sqlc.yaml`/lint (`golangci-lint`) configuration yet — don't assume they exist.

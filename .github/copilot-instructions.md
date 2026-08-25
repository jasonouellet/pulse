# PULSE — Instructions pour agents IA

Plateforme SaaS de gestion de sports collectifs (multi-sport, multi-niveaux, rosters éphémères, tournois, calendrier dynamique). Monolithe modulaire en Go conçu pour une extraction future vers des microservices.

## Architecture

Voir `docs/ARCHITECTURE.md`, `docs/ArchitectureTechnique.md`, `docs/ArchitectureC4.md` et les ADRs dans `docs/adr/` (source de vérité pour tout choix structurant). Toujours valider un changement contre ces documents avant de le faire.

- **Pattern** : Hexagonal (Ports & Adapters) par module, à l'intérieur d'un monolithe modulaire (`cmd/monolith/main.go`).
- **Modules** (`internal/{module}/`) : `core`, `tournament`, `scheduling`, `finance`, `evaluation`. Chacun suit la structure :
  ```
  internal/{module}/
  ├── domain/    # logique métier pure, sans dépendance externe
  ├── ports/     # interfaces (contrats d'entrée/sortie)
  └── adapters/postgres/queries/{entity}.sql  # requêtes façon sqlc
  ```
- **Isolation stricte des modules** : toute communication inter-module passe par une interface Go (port). **Aucun JOIN SQL entre schémas PostgreSQL n'est autorisé** — un schéma Postgres par module (`core`, `tournament`, `scheduling`, `finance`, `evaluation`).
- **Multi-sport** : toutes les entités pertinentes portent un `sport_id` (UUID) pour préparer le support multi-sport (Soccer en phase 1).

## Stack technique

- **Backend** : Go 1.22, routeur `go-chi/chi/v5`, driver `jackc/pgx/v5`, logging structuré via `slog` (JSON handler).
- **Base de données** : PostgreSQL 16+ (schéma par module), Redis 7+ (cache/pub-sub/verrous).
- **Frontend** : React + TypeScript, Tailwind CSS, Radix UI/Shadcn — dossier `frontend/` encore vide (scaffold à venir).
- **Conteneurs** : builds Docker multi-stage, images distroless non-root pour le backend.

## Conventions

- **Migrations SQL** : `migrations/{module}/NNNNNN_description.up.sql` (convention `golang-migrate`), une migration par schéma/module.
- **Requêtes SQL** : fichiers `internal/{module}/adapters/postgres/queries/{entity}.sql`, commentaires style sqlc (`-- name: NomFonction :one|:many|:exec`).
- **Schéma DB** : clés primaires en UUID (`gen_random_uuid()`), `TIMESTAMP WITH TIME ZONE`, documentation via `COMMENT ON TABLE|COLUMN`.
- **Langue** : tout le contenu (code, identifiants, commentaires, messages de commit, PRs) doit être rédigé en anglais, à l'exception du répertoire `docs/` qui reste en français.
- **Accessibilité (frontend)** : WCAG 2.1 AA, mobile-first, cibles tactiles min. 44x44px, support natif Dark/Light via variables CSS.
- Après toute tâche terminée, rappeler à l'utilisateur de mettre à jour `docs/ROADMAP_TODO.md`.

## Build & exécution

```bash
cp .env.example .env
docker compose up -d          # backend + PostgreSQL + Redis (+ frontend si présent)
curl http://localhost:8080/healthz
```

Il n'y a pas encore de Makefile, ni de suite de tests, ni de configuration `sqlc.yaml`/lint (`golangci-lint`) — ne pas supposer leur existence.

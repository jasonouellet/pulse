# Project PULSE (PULSE OS)

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-18.x-61DAFB?style=flat&logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16%2B-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-Source--Available_Non--Commercial-blue.svg)](docs/adr/ADR-005.md)
[![WCAG 2.1 AA](https://img.shields.io/badge/Accessibility-WCAG_2.1_AA-success.svg)](https://www.w3.org/WAI/standards-guidelines/wcag/)

**PULSE OS** est un système de gestion moderne, modulaire et hautement disponible conçu pour les clubs sportifs jeunesse multi-sports (avec un focus initial sur le soccer). Il prend en charge l'administration des bassins d'âges, la constitution d'équipes éphémères, la sous-inscription à la carte pour les événements, ainsi que la logistique dynamique des terrains et des tournois.

---

## 🏗️ Architecture & Stack Technique

Le projet est conçu selon les principes de la **Clean Architecture (Ports & Adapters)** sous forme de **Monolithe Modulaire Go**, garantissant une séparation stricte des domaines métier et permettant une migration vers des microservices sans refactorisation lourde (`ADR-003`).

### Backend
* **Langage :** Go 1.22+
* **Routeur HTTP :** Chi Router (API RESTful)
* **Base de données :** PostgreSQL 16+ avec isolation stricte par schémas (`core`, `tournament`, `scheduling`, `finance`)
* **Génération SQL :** `sqlc` (requêtes type-safe compilées en Go)
* **Migrations :** `golang-migrate` (exécution automatique au démarrage)
* **Observabilité :** OpenTelemetry OTLP/gRPC (`otelchi`, `otelslog`, `otelpgx`) (`ADR-006`)

### Frontend
* **Framework :** React + Vite (TypeScript)
* **Styling & UI :** Tailwind CSS, Radix UI (Thème Day/Dark, accessibilité WCAG 2.1 AA)
* **Icônes :** Lucide React
* **Testing :** Vitest + React Testing Library + Happy-DOM + Vitest-Axe
* **Observabilité Web :** OpenTelemetry Web SDK (`@opentelemetry/sdk-trace-web`)

---

## 📁 Structure du Projet

```text
pulse/
├── cmd/
│   └── monolith/          # Point d'entrée du serveur Go Monolithe
├── internal/
│   ├── core/              # Module Core (Users, Pools, Sports abstraction)
│   ├── tournament/        # Module Tournois (Rosters éphémères, Brackets)
│   ├── scheduling/        # Module Horaires (Fields, Matches, Attendance)
│   └── finance/           # Module Finance (Expenses, Payments)
├── migrations/            # Scripts SQL par schéma (000001_create_*.up.sql)
├── pkg/
│   ├── database/          # Connexion PGX & Runner de migrations
│   └── observability/     # Pipeline OpenTelemetry Tracer & Loggers
├── frontend/              # Application Vite + React + TypeScript
└── docs/
    ├── adr/               # Architecture Decision Records (001 à 006)
    └── ROADMAP_TODO.md    # Suivi d'avancement du projet

## 🚀 Démarrage Rapide (Développement Local)

### Prérequis
* **Go** 1.22 ou plus récent
* **Node.js** v18+ & `npm`
* **PostgreSQL** 16+
* **Podman** / **Docker** (Optionnel pour la stack complète)

### 1. Démarrer le Backend Go

Définissez vos variables d'environnement et lancez le serveur :

```bash
# Variables d'environnement de dev
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres_secret
export DB_NAME=pulse_db
export DB_SSLMODE=disable
export PORT=8080

# Lancer le monolithe
go run cmd/monolith/main.go


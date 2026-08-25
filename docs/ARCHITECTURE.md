# ARCHITECTURE TECHNIQUE — PROJECT PULSE

## 1. Vue d'Ensemble C4

- **Frontend :** SPA React.js avec TypeScript, Tailwind CSS et Radix UI (PWA & Responsive, Day/Dark Mode).
- **Backend :** Go (Golang) utilisant le pattern Monolithe Modulaire (Modular Monolith).
- **Stockage :** PostgreSQL (Isolation par schémas SQL : `core`, `tournament`, `scheduling`, `finance`, `evaluation`).
- **Cache & Real-time :** Redis (Cache HTTP, verrous distribués et Pub/Sub WebSockets pour les scores).

## 2. Abstraction Multi-Sports (`ADR-004`)

Toutes les entités clés possèdent une référence `sport_id` (définie dans `core.sports`). En Phase 1, la valeur par défaut est associée au Soccer.

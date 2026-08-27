# ARCHITECTURE TECHNIQUE

## 1. Vue d'Ensemble C4

* **Frontend :** SPA React.js avec TypeScript, Tailwind CSS et Radix UI (PWA & Responsive, Day/Dark Mode).
* **Backend :** Go (Golang) utilisant le pattern Monolithe Modulaire (Modular Monolith).
* **Stockage :** PostgreSQL (Isolation par schémas SQL : `core`, `tournament`, `scheduling`, `finance`, `evaluation`).
* **Cache & Real-time :** Redis (Cache HTTP, verrous distribués et Pub/Sub WebSockets pour les scores).

## 2. Abstraction Multi-Sports (`ADR-004`)

Toutes les entités clés possèdent une référence `sport_id` (définie dans `core.sports`).
En Phase 1, la valeur par défaut est associée au Soccer.

## 2. Diagramme de Conteneurs & Modules (Go)

```text
/cmd/monolith/main.go   --> Lance tous les modules dans 1 binaire
/internal/
  /core                 --> Schéma SQL: core
  /tournament           --> Schéma SQL: tournament
  /scheduling           --> Schéma SQL: scheduling
  /finance              --> Schéma SQL: finance
  /evaluation           --> Schéma SQL: evaluation
```

## 3. Modèle d'Isolation (Ports & Adapters)

Les modules communiquent uniquement via des interfaces Go (internal/module/ports).

En mode Monolithe, l'adaptateur est en mémoire Go.

En mode Microservice (futur), l'adaptateur bascule en gRPC sans modifier la logique métier.

### 4. Registre ADR

* **`ADR-001-frontend-stack.md` :** Choix de React, TypeScript, Tailwind, Radix UI et norme WCAG 2.1 AA.
* **`ADR-002-backend-and-data-stack.md` :** Choix de Go, PostgreSQL (sqlc) et Redis.
* **`ADR-003-modular-monolith-to-microservices.md` :** Choix du Monolithe Modulaire avec migration microservices sans refactoring via Ports & Adapters.
*

---

## [TODO ARCHITECTURE]

[ ] Structurer la gateway WebSocket en Go pour pousser les scores depuis Redis.
[ ] Valider la stratégie de fallback hors-ligne (Offline PWA) pour la saisie des présences par les entraîneurs sur le terrain.

## 5. Fichier `docs/ROADMAP_TODO.md` (Roadmap & Suivi des tâches)

## 💡 Reste à Réfléchir / À Clarifier

1. **Paiements :** Stripe Connect pour ventiler les paiements directement aux clubs ou compte marchand unique ?
2. **Notifications :** Préférer les notifications web push ou l'intégration Twilio SMS pour les urgences du terrain (ex: annulation pour orage) ?

# ROADMAP ET SUIVI DES TÂCHES — PROJECT PULSE

**Nom de code projet :** Project PULSE
**Architecture :** Monolithe Modulaire Go (Clean Architecture / Ports & Adapters), PostgreSQL 16+ (Isolation par schémas SQL), Redis, React.js (TypeScript, Tailwind, Radix UI).
**Licence & Propriété :** Source-Available Non-Commercial (ADR-005, CLA.md) — Gratuit pour auto-hébergement des clubs à but non lucratif.
**Périmètre :** Multi-sports collectifs (Focus Phase 1 : Soccer).

---

## 📊 Tableau de Bord de l'Avancement (Phase 1)

| Domaine / Module                                            |    Statut    | Progression |
| :---------------------------------------------------------- | :----------: | :---------: |
| **Ingénierie & Context Engineering**                        | **TERMINÉ**  |    100%     |
| **Architecture C4 & ADRs (001 à 006 + CLA)**                | **EN COURS** |     90%     |
| **Analyse Comparative & Validation des Fichiers Tournois**  | **TERMINÉ**  |    100%     |
| **Backend Go & Architecture Hexagonale (Core Module)**      | **EN COURS** |     85%     |
| **Schéma de Données SQL (`core` révisé, `tournament`, etc.) | **EN COURS** |     30%     |
| **Environnement DevContainer & Observabilité OTEL**         | **EN COURS** |     90%     |
| **Frontend React (Setup Vite, Tailwind & Theme Day/Dark)**  | **EN COURS** |     40%     |
| **Gouvernance de dépôt (Licence, CONTRIBUTING, SECURITY)**  | **EN COURS** |     60%     |

---

## 🚀 PHASE 1 : Fondations Techniques & Data Model (EN COURS)

### 1.1 — Context Engineering, Architecture & Licence (✅ Livré)

* [x] Définir les principes d'architecture et le découpage fonctionnel des 3 vues (Admin, Parent, Sportif).
* [x] Concevoir la stratégie des Bassins d'âges (Pools) et Rosters Éphémères.
* [x] Définir le modèle de sous-inscriptions à la carte (Tournois #1, #3, #4).
* [x] Définir le moteur d'organisation d'événements (découpage terrains 11v11 ➔ 7v7, 3 matchs min. garantis).
* [x] Rédiger le registre d'ADRs complet :
  * [x] `ADR-001` : Frontend React, TypeScript, Tailwind, Radix UI, WCAG 2.1 AA.
  * [x] `ADR-002-B` : Backend Go compilé, PostgreSQL 16+ et cache Redis.
  * [x] `ADR-003` : Monolithe Modulaire avec migration microservices sans refactoring (Ports & Adapters).
  * [x] `ADR-004` : Abstraction Multi-Sports (`sport_id`) avec focus initial Soccer.
  * [x] `ADR-005` : Modèle de Licence Source-Available Non-Commercial (Auto-hébergement gratuit pour les clubs).
  * [x] `ADR-006` : Stratégie d'Observabilité globale via OpenTelemetry (Traces, Metrics, Logs pour Go, React, Postgres, Redis).
* [x] Rédiger l'accord de licence de contributeur (`CLA.md`).
* [x] Générer la documentation d'architecture C4 complète (Diagrammes Mermaidjs Niveaux 1 à 3).
* [x] Établir le comparatif de marché et la validation légale des marques.
* [x] Valider l'adéquation de l'architecture avec les fichiers réels de tournois (_Rimouski Juin/Juillet/Août 2026_).
* [x] Fixer le nom de code projet : **Project PULSE**.

---

### 1.1.B — Dette technique identifiée à l'audit du 28 août 2026 (🔴 À résorber avant de poursuivre)

Un audit complet du dépôt sur `feature/core` a révélé plusieurs écarts entre ce que ce ROADMAP déclarait "fait" et l'état réel du code. Cette section documente ce qui doit être réglé **avant** d'attaquer les schémas `tournament`/`scheduling`/`finance`, pour éviter d'empiler de la dette sur de la dette.

#### Migrations & données

* [ ] Écrire `migrations/core/000001_create_core_schema.down.sql` — actuellement absent malgré la case cochée ; aucun rollback possible en l'état.
* [ ] Créer `sqlc.yaml` et configurer sqlc, ou retirer `internal/core/adapters/postgres/queries/*.sql` tant que ce n'est pas prêt — ces fichiers ne génèrent rien actuellement.
* [ ] Résoudre le conflit : `CreatePool` est défini deux fois (`users.sql` et `pools.sql`) avec des signatures incompatibles. sqlc refusera de générer tant que ce n'est pas réglé.
* [ ] Décider si `internal/core/adapters/postgres/user_repository.go` (SQL à la main) migre vers sqlc une fois configuré, ou si les deux approches coexistent délibérément — actuellement c'est un flottement non documenté.

#### Code mort & duplication

* [ ] Supprimer `internal/pkg/telemetry/tracer.go` (implémentation OTEL stdout orpheline, jamais appelée — `pkg/observability/otel.go` est la seule utilisée par `main.go`) ou fusionner s'il y a une intention derrière (ex. mode debug local).

#### Architecture Decision Records

* [ ] Renuméroter l'un des deux `ADR-006` (`ADR-006-observability-opentelemetry.md` vs `ADR-006-telemetry-and-health-obligations.md`) — collision de numérotation, contenus qui se chevauchent partiellement.
* [ ] Rédiger un ADR pour l'adoption de `huma` (routing + génération OpenAPI), qui remplace l'approche `chi` brut + `swaggo/swag` initialement documentée — changement structurant non tracé.
* [ ] Mettre à jour `docs/ARCHITECTURE.md` pour refléter `huma` (encore rédigé autour de `swaggo/swag`).

#### Outillage & CI

* [ ] Trancher entre `dependabot.yml` et `renovate.json` — les deux sont configurés sur les mêmes écosystèmes (Go, npm, Docker, GitHub Actions), ce qui va générer des PRs en double. Garder un seul outil.
* [ ] Corriger `renovate.json` : la clé `matchMatchers` n'existe pas dans le schéma Renovate (c'est `matchManagers`) — le groupement par écosystème ne fonctionne pas tel quel.

#### Question de priorisation (pas un bug)

* [ ] Statuer sur l'ordre d'investissement : des manifests Kubernetes complets (`deployments/kubernetes/kind/`) et un `otel-collector` existent déjà, alors que les schémas `tournament`, `scheduling`, `finance`, `evaluation` n'ont pas démarré et que Redis n'est référencé nulle part dans le code Go. À valider consciemment plutôt que de laisser l'infra continuer à devancer le domaine métier.

---

* [x] Écrire `migrations/core/000001_create_core_schema.up.sql` (Anglais strict, UUID `gen_random_uuid()`) — voir 1.1.B pour le `.down.sql` manquant
  * [x] Table `core.sports` (Multi-sport abstraction via JSONB rules).
  * [x] Tables `core.users`, `core.parents_children`, `core.player_profiles`.
  * [x] Table `core.pools` (Bassins d'âges / Catégories).
  * [x] Ajouter les instructions `COMMENT ON` pour la documentation automatique.
* [x] Écrire les requêtes préparées `sqlc` (`internal/core/adapters/postgres/queries/users.sql`).
* [ ] Écrire `migrations/tournament/000001_create_tournament_schema.up.sql`
  * [ ] Tables `tournament.rosters`, `tournament.event_sub_registrations`, `tournament.brackets`.
* [ ] Écrire `migrations/scheduling/000001_create_scheduling_schema.up.sql`
  * [ ] Tables `scheduling.fields`, `scheduling.sub_fields`, `scheduling.matches`, `scheduling.attendances`.
* [ ] Écrire `migrations/finance/000001_create_finance_schema.up.sql`
  * [ ] Tables `finance.event_expenses`, `finance.referee_payments`.
* [ ] Configurer `sqlc.yaml` pour l'isolation stricte des requêtes SQL par module Go.

---

### 1.3 — Socle Backend Go (🔄 En cours)

* [x] Mettre en place la structure Hexagonale du dépôt Go (`/cmd/backend`, `/internal/core`, etc.).
* [x] Implémenter le Port `UserRepository` et l'Adaptateur PostgreSQL pour le module Core.
* [x] Implémenter le contrôleur HTTP pour l'API `/api/v1/core/users`, migré vers `huma` (génération OpenAPI automatique, validation de schéma).
* [x] Exposer les sondes `/livez`, `/readyz`, `/healthz` et les métriques Prometheus (`/metrics`).
* [x] Configurer le serveur HTTP Chi avec Graceful Shutdown dans `cmd/backend/main.go`.
* [x] Mettre en place la suite de tests unitaires HTTP avec Mocks (`internal/core/adapters/http/user_handler_test.go`), migrée vers `huma`.
* [x] Intégrer `golang-migrate` pour l'exécution automatique des migrations SQL au démarrage.
* [x] Configurer l'exportation OpenTelemetry (OTLP/gRPC) dans le backend (`pkg/observability/otel.go`).
* [ ] Nettoyer le doublon `internal/pkg/telemetry/tracer.go` (voir 1.1.B).
* [ ] Augmenter la couverture de tests.

---

### 1.4 — Socle Frontend React (🔄 En cours)

* [x] Initialiser le projet Vite + React + TypeScript dans `frontend/`.
* [x] Configurer Tailwind CSS, Radix UI et le système de thèmes (Dark/Day mode).
* [x] Intégrer les packages d'instrumentation web OpenTelemetry (`@opentelemetry/sdk-trace-web`).
* [x] Créer les composants de base accessibles WCAG (Button, Modal, Dynamic Data Table, Card).
* [x] Mettre en place le dictionnaire de termes par sport (Lexique Soccer).
* [x] Configurer l'exportation OpenTelemetry.
* [ ] Augmenter la couverture de tests.

---

## Amélioration continue SLDC (🔄 En cours)

### Z — Tooling

* [x] Configurer la stack d'environnement DevContainer (`.devcontainer/devcontainer.json`, `Dockerfile`).
* [x] Scaffolder les manifests Kubernetes (Kind) — voir 1.1.B pour la question de priorisation.
* [ ] Corriger `dependabot.yml` / trancher entre Dependabot et Renovate — voir 1.1.B.
* [ ] Ajouter les outils [Score](https://docs.score.dev/docs/score-implementation/).

### Z.2 — Gouvernance de dépôt

* [x] Rédiger `README.md` (anglais).
* [ ] Pousser `CONTRIBUTING.md` (rédigé, anglais — voir 1.1.B).
* [ ] Pousser `SECURITY.md` (rédigé, anglais — voir 1.1.B).
* [ ] Corriger `LICENSE` (coquille + date de conversion — voir 1.1.B).

---

## 🔮 PHASES SUIVANTES (Aperçu)

### PHASE 2 : Core Domain & Portail Parent

* [ ] Authentification & Gestion Granulaire des Rôles (RBAC).
* [ ] Espace Parent : Fiche Enfant, Inscription au Bassin U9-U10.
* [ ] Tunnel de Sous-Inscriptions aux tournois optionnels.

### PHASE 3 : Équipes Éphémères & Calendrier Dynamique

* [ ] Outil d'administration : Répartition des jeunes (Draft / Team Builder).
* [ ] Moteur API de consolidation du calendrier familial (Pratiques du bassin + Matchs du roster).
* [ ] Module RSVP Express (Confirmation de présence 1-click).

### PHASE 4 : Logistique Terrains, Brackets & Finances

* [ ] Algorithme d'affectation et de découpage des terrains (11v11 ➔ 2x 9v9 / 4x 7v7).
* [ ] Générateur d'arbres de tournois avec repêchage (Consolidation 3 matchs min. garantis).
* [ ] Dashboard financier par événement (Postes de dépenses vs Sous-inscriptions).

### PHASE 5 : Évaluations & Real-Time

* [ ] Grilles d'évaluation de compétences par calibre (Technique, Tactique, Physique, Mental).
* [ ] Gateway WebSocket Go + Redis Pub/Sub pour la diffusion des scores en direct sur le terrain.

---

## 💡 Reste à Réfléchir / Décisions à Prendre

1. **Paiements :** Choix entre Stripe Connect (ventilation automatique par club) ou Stripe Checkout standard.
2. **Notifications Urgentes :** Priorité aux WebPush ou intégration Twilio SMS pour les annulations de terrain sur le coup (ex: orages).
3. **Identité Finale :** Choix définitif de la marque commerciale avant la mise en production.
4. **Évolutivité :** Préparer le projet pour une future migration vers une architecture microservices si nécessaire (`ADR-003`).
5. **Multi-Sports :** Définir la stratégie d'ajout de nouveaux sports (soccer, basketball, volleyball) et l'impact sur le modèle JSONB (`ADR-004`).
6. **SDLC & CI/CD :** Définir les pipelines GitHub Actions pour la validation des PRs, les tests automatiques et les builds d'images distroless OCI.
7. Comment nommer le rôle "Parent" pouur qu'il représente autant un tuteur, un grand parent ou un grand frère/soeurs
8.

# ROADMAP ET SUIVI DES TÂCHES — PROJECT PULSE

**Nom de code projet :** Project PULSE
**Architecture :** Monolithe Modulaire Go (Clean Architecture / Ports & Adapters / Huma v2 / Chi), PostgreSQL 16+ (Isolation par schémas SQL - ADR-003), Redis, React.js (TypeScript, Tailwind, Radix UI).
**Licence & Propriété :** Source-Available Non-Commercial (ADR-005, CLA.md) — Auto-hébergement gratuit pour les clubs à but non lucratif.
**Périmètre :** Multi-sports collectifs (Focus Phase 1 : Soccer mineur, catégories U9-U10F à U18).

---

## 📊 Tableau de Bord de l'Avancement (Phase 1)

| Domaine / Module                                                                     |    Statut    | Progression |
| :----------------------------------------------------------------------------------- | :----------: | :---------: |
| **Ingénierie & Context Engineering & Cadrage Fonctionnel**                           | **TERMINÉ**  |    100%     |
| **Architecture C4 & ADRs (001 à 007 + CLA)**                                         | **TERMINÉ**  |    100%     |
| **Backend Go & Architecture Hexagonale (Core Module & Huma)**                        | **EN COURS** |     85%     |
| **Schémas de Données SQL (`core` révisé, `tournament`, `scheduling`, `evaluation`)** | **EN COURS** |     50%     |
| **Environnement DevContainer, Tooling CI & Observabilité (Prometheus/OTEL)**         | **EN COURS** |     90%     |
| **Frontend React (Setup Vite, Tailwind & Theme Day/Dark, WCAG)**                     | **EN COURS** |     40%     |
| **Gouvernance de dépôt (Licence, CONTRIBUTING, SECURITY)**                           | **EN COURS** |     70%     |

---

## 📌 Décision Terminologique : Rôle "Guardian" (Représentant Familial)

> **Décision validée, mise à jour 2026-08-29 pour refléter l'implémentation réelle** (la version précédente de cette note décrivait un plan qui a dérivé pendant la migration `core`) :
>
> * **En base de données (`core`) :** rôle `GUARDIAN` dans `core.user_role` ; table de jonction `core.parents_children` avec `relationship_type` (Enum: `MOTHER`, `FATHER`, `LEGAL_GUARDIAN`, `OTHER`).
> * **Dans l'interface UX/UI :** libellé affiché **"Représentant"** / **"Mes enfants"**.
> * **[ ] À trancher :** `GRANDPARENT`/`SIBLING` avaient été envisagés dans la version précédente de cette décision — à ajouter à `relationship_type` si un club en a besoin, ou à laisser dans `OTHER` pour l'instant.

---

## 🚀 PHASE 1 : Fondations Techniques & Data Model (EN COURS)

### 1.1 — Context Engineering, Architecture, Cadrage & Licence (✅ Livré)

* [x] Définir les principes d'architecture et le découpage fonctionnel des 3 vues (Admin/DT, Staff/Entraîneur, Représentant Familial & Sportif).
* [x] Établir le **Dossier de Cadrage Fonctionnel (BIZBOK® / BABOK® / Strategyzer / VSM)** (`docs/CADRAGE_FONCTIONNEL.md`) minimal pour le compléter itérativement.
* [x] Rédiger le registre d'ADRs complet (à date) :
  * [x] `ADR-001` : Frontend React, TypeScript, Tailwind, Radix UI, WCAG 2.1 AA.
  * [x] `ADR-002-B` : Backend Go compilé, PostgreSQL 16+ et cache Redis.
  * [x] `ADR-003` : Monolithe Modulaire avec migration microservices sans refactoring (Ports & Adapters).
  * [x] `ADR-004` : Abstraction Multi-Sports (`sport_id`) avec focus initial Soccer.
  * [x] `ADR-005` : Modèle de Licence Source-Available Non-Commercial.
  * [x] `ADR-006` : Stratégie d'Observabilité globale (Prometheus `/metrics` + OpenTelemetry).
  * [x] `ADR-007` : Adéquation du framework **Huma v2** pour la couche API OpenAPI/Swagger UI.
* [x] Rédiger l'accord de licence de contributeur (`CLA.md`).
* [x] Générer la documentation d'architecture C4 complète (Diagrammes Niveaux 1 à 3).
* [x] Établir le comparatif de marché et la validation légale des marques.
* [x] Valider l'adéquation de l'architecture avec les fichiers réels de tournois (*Rimouski Juin/Juillet/Août 2026*).
* [x] Fixer le nom de code projet : **Project PULSE**.

---

### 1.1.B — Résorption de la Dette Technique & Qualité CI (🔄 En cours)

#### Migrations & données

* [ ] Écrire `migrations/core/000001_create_core_schema.down.sql` pour permettre les rollbacks.
* [ ] Configurer `sqlc.yaml` à la racine pour cibler l'isolation par schéma.
* [ ] Résoudre le conflit de double définition `CreatePool` entre `users.sql` et `pools.sql`.
* [ ] Aligner l'adaptateur `user_repository.go` sur les requêtes générées par `sqlc`.

#### Code mort & observabilité

* [x] Nettoyer le doublon d'initialisation OTEL : suppression de `internal/pkg/telemetry/tracer.go` au profit du package officiel `pulse/pkg/observability`.
* [x] Exposer le handler `/metrics` Prometheus (`prometheus/promhttp`) sur le routeur Chi central dans `cmd/backend/main.go`.

#### Outillage & CI/CD

* [x] Corriger le tag d'image de base Go dans `deployments/docker/backend.Dockerfile` (`golang:1.23-alpine`).
* [x] Régler les erreurs du linter Markdown `markdownlint-cli2`.
* [x] Intégrer le hook `markdownlint-cli2` dans `.pre-commit-config.yaml` et corriger l'avertissement de chemin Hadolint `(.*/)?Dockerfile.*`.
* [x] Résoudre les erreurs de formatage SQLFluff sur les migrations (respect de la limite de longueur de ligne de 160 caractères).
* [ ] Trancher entre `dependabot.yml` et `renovate.json` (supprimer l'un des deux pour éviter les PRs en double).
  * dependabot, c'est pour la gestion active, du repo, renovate c'est pour updater localement ou sur demande.
* [x] Corriger la clé `matchMatchers` ➔ `matchManagers` dans `renovate.json`.
* [ ] Évaluer l'intégration du Modèle Score (`score-compose`) pour la génération des manifests de dev local.

---

### 1.2 — Modélisation des Données & Schémas SQL (🔄 En cours)

* [x] Écrire `migrations/core/000001_create_core_schema.up.sql` (`core.sports`, `core.users`, `core.parents_children`, `core.player_profiles`, `core.pools`).
* [x] Écrire `migrations/tournament/000001_create_tournament_schema.up.sql` — restructuré selon `ADR-008`, voir 1.2.B.
* [x] Écrire `migrations/scheduling/000001_create_scheduling_schema.up.sql` (`date_windows`, `fields`, `sub_fields`, `activities`, `attendances`).
* [ ] Écrire `migrations/evaluation/000001_create_evaluation_schema.up.sql` (`player_ratings`).
* [ ] Réserver la structure `migrations/finance/` pour la phase v2.0 (Out-of-Scope v1.0).

#### 1.2.B — Restructuration selon `ADR-008` (✅ Migrations faites — reste le câblage applicatif)

* [x] Déplacer `TRAINING_GROUP`/`SEASON_TEAM` de `tournament.roster_type` vers `core` (`core.teams`, `core.team_pools`, `core.team_players`).
* [x] `tournament.roster_type` retiré — un seul type restait (`EVENT_ROSTER`), l'enum n'avait plus de raison d'être.
* [x] Ajouter `tournament.event_eligibility` (une ou plusieurs plages âge/genre par événement).
* [x] Ajouter `tournament.team_entries` (club + événement + bassin source) — distinct de l'alignement (`tournament.rosters`, maintenant 1:1 avec `team_entries`).
* [x] Créer `migrations/scheduling/000001_create_scheduling_schema.up.sql` comme second noyau partagé, sans dépendance sortante.
* [x] Retirer `start_date`/`end_date` de `core.pools` et `tournament.events`, référencer `scheduling.date_windows` à la place (`window_id`).
* [x] Ajouter `parent_club_id` (auto-référencé) et `org_type` (`CLUB`/`ASSOCIATION`/`FEDERATION`) sur `core.clubs`.
* [ ] Renommer `roster` → `alignement` dans le code/UI en français ; garder `roster` en anglais/code — **pas encore fait côté frontend** (`RosterBuilderPage`, `NAV_BY_ROLE`, etc. disent encore "roster"/"Rosters").
* [ ] **🆕 Bloquant :** `pkg/database/migrate.go`/`main.go` doivent appeler les migrations dans l'ordre `scheduling` → `core` → `tournament`, avec une table de suivi (`x-migrations-table`) distincte par schéma — sinon `golang-migrate` croit que tout est déjà appliqué après le premier appel.
* [ ] Écrire les Ports/Repositories Go pour `core.clubs` (hiérarchie), `core.teams`, `scheduling.*`, `tournament.*` — aucun n'existe encore.
* [ ] Semer des données de base (`scheduling.date_windows` pour la saison en cours, positions par sport) — les tables existent mais sont vides.

---

### 1.3 — Socle Backend Go & API Huma v2 (🔄 En cours)

* [x] Mettre en place la structure Hexagonale du dépôt Go.
* [x] Implémenter le Port `UserRepository` et l'Adaptateur PostgreSQL pour le module Core.
* [x] Configurer le routeur Chi avec Huma v2 (`/docs` Swagger UI, `/openapi` spec JSON).
* [x] Exposer les sondes `/livez`, `/readyz`, `/healthz` et `/metrics` Prometheus.
* [x] Mettre en place le Graceful Shutdown et l'exécution automatique des migrations via `golang-migrate`.
* [x] Implémenter les handlers Huma pour `core` (`/api/v1/core/users`).
* [ ] Implémenter les handlers Huma pour `tournament`.

---

### 1.4 — Socle Frontend React (🔄 En cours)

* [x] Initialiser le projet Vite + React + TypeScript dans `frontend/`.
* [x] Configurer Tailwind CSS, Radix UI et le système de thèmes.
* [x] Créer les composants de base accessibles WCAG 2.1 AA.
* [x] Intégrer l'instrumentation web OpenTelemetry.
* [ ] Implémenter le composant de sélection des joueurs par Drag & Drop.
* [ ] Créer la vue d'affichage du calendrier unifié.

## 2. Conformité Légale & Sécurité (Loi 25 Québec / PIPEDA Canada)

### 2.1 Consentements & Gestion des Données Mineurs (< 14 ans)

* [ ] **[Loi 25]** Implémenter le flux d'onboarding avec opt-in explicite du *Family Guardian* pour la collecte de données des mineurs de moins de 14 ans.
* [ ] **[Loi 25]** Créer un consentement séparé et explicite pour le traitement des notes vocales et l'analyse sportive par IA.
* [ ] **[Loi 25]** Rédiger la politique de confidentialité en langage clair et afficher les coordonnées du Responsable de la protection des données.

### 2.2 Droits des Utilisateurs & Anonymisation IA

* [ ] **[Loi 25]** Développer l'endpoint de droit à l'oubli (`DELETE /api/v1/core/me`) et d'exportation structurée des données (portabilité JSON/CSV).
* [ ] **[Loi 25]** Implémenter la couche d'anonymisation dans l'adaptateur IA Go (masquage des noms/prénoms avant envoi aux LLMs tiers comme Groq/OpenAI).
* [ ] **[Loi 25]** Configurer un job de purge/archivage automatique des données pour les comptes inactifs.

## 🔮 PHASES SUIVANTES & FEUILLE DE ROUTE PRODUIT

### PHASE 2 : Gestion d'Équipes Avancée, Horaires & Évaluation Light (Priorité Cible)

> Ordre de priorité confirmé : **1) saisons & entraînements → 2) tournois locaux → 3) tournois régionaux/provinciaux.** Chaque étape s'appuie sur la précédente ; ne pas paralléliser au-delà de ce qui est nécessaire.

* [ ] **Priorité 1 — Saisons & Entraînements :** Finalisation de `core` (bassins, groupes d'entraînement, équipes de saison — voir 1.2.B) et du calendrier de pratiques (`scheduling`, EPIC-2).
* [ ] **Priorité 2 — Tournois Locaux :** `EPIC-1` (éligibilité, équipes engagées, alignements) et `EPIC-2` (terrains, calendrier de matchs) pour des événements à un seul club/une seule région.
* [ ] **Priorité 3 — Tournois Régionaux/Provinciaux :** Extension multi-clubs de `EPIC-1` (plusieurs clubs engagent des équipes dans un même événement) — dépend de la hiérarchie organisationnelle (`ADR-008 §6`) si la fédération/association doit superviser.
* [ ] **EPIC-3 :** Saisie des fiches d'évaluation simplifiées par l'entraîneur.
* [ ] **Espace Représentant Familial :** Consultation des équipes, des calendriers et confirmation de présence.

---

### PHASE 3 : Portail Financier & Interactivité Temps Réel (v2.0 — Out-of-Scope v1.0)

* [ ] **EPIC-4 (Finance) :** Module de facturation et sous-inscriptions aux tournois via Stripe.
* [ ] **EPIC-5 (Scheduling Extended) :** RSVP Express 1-click et notifications Push/SMS d'urgence.
* [ ] **EPIC-6 (Evaluation Extended) :** Grilles d'évaluation personnalisables avec schémas JSONB.

---

### PHASE 4 : Logistique Complexe & Moteur de Tournoi (v3.0+)

* [ ] **EPIC-7 (Tournament Extended) :** Générateur automatique d'arbres de tournoi et diffusion des scores en direct via WebSockets + Redis.
* [ ] **EPIC-8 (Arbitrage) :** Module de désignation des trios d'arbitres et calcul des paies.
* [ ] **PWA Offline-First :** Prise de présence et saisie des évaluations sur le terrain sans réseau.

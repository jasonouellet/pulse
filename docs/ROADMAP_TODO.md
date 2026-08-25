# ROADMAP ET SUIVI DES TÂCHES — PROJECT PULSE

**Nom de code projet :** Project PULSE (PULSE OS)  
**Architecture :** Monolithe Modulaire Go (Clean Architecture / Ports & Adapters), PostgreSQL 16+ (Isolation par schémas SQL), Redis, React.js (TypeScript, Tailwind, Radix UI).
**Licence & Propriété :** Source-Available Non-Commercial (ADR-005, CLA.md) — Gratuit pour auto-hébergement des clubs à but non lucratif.
**Périmètre :** Multi-sports collectifs (Focus Phase 1 : Soccer).

---

## 📊 Tableau de Bord de l'Avancement (Phase 1)

| Domaine / Module | Statut | Progression |
| :--- | :---: | :---: |
| **Ingénierie & Context Engineering** | **TERMINÉ** | 100% |
| **Architecture C4 & ADRs (001 à 005 + CLA)** | **TERMINÉ** | 100% |
| **Analyse Comparative & Modèle de Noms** | **TERMINÉ** | 100% |
| **Schéma de Données SQL (`core`, `tournament`, etc.)** | **À FAIRE** | 0% |
| **Backend Go (Scaffolding Monolithe & sqlc)** | **À FAIRE** | 0% |
| **Frontend React (Setup Vite, Tailwind & Theme Day/Dark)** | **À FAIRE** | 0% |

---

## 🚀 PHASE 1 : Fondations Techniques & Data Model (EN COURS)

### 1.1 — Context Engineering, Architecture & Licence (✅ Livré)

- [x] Définir les principes d'architecture et le découpage fonctionnel des 3 vues (Admin, Parent, Sportif).
- [x] Concevoir la stratégie des Bassins d'âges (Pools) et Rosters Éphémères.
- [x] Définir le modèle de sous-inscriptions à la carte (Tournois #1, #3, #4).
- [x] Définir le moteur d'organisation d'événements (découpage terrains 11v11 ➔ 7v7, 3 matchs min. garantis).
- [x] Rédiger le registre d'ADRs complet :
  - [x] `ADR-001` : Frontend React, TypeScript, Tailwind, Radix UI, WCAG 2.1 AA.
  - [x] `ADR-002-B` : Backend Go compilé, PostgreSQL 16+ et cache Redis.
  - [x] `ADR-003` : Monolithe Modulaire avec migration microservices sans refactoring (Ports & Adapters).
  - [x] `ADR-004` : Abstraction Multi-Sports (`sport_id`) avec focus initial Soccer.
  - [x] `ADR-005` : Modèle de Licence Source-Available Non-Commercial (Auto-hébergement gratuit pour les clubs).
- [x] Rédiger l'accord de licence de contributeur (`CLA.md`).
- [x] Générer la documentation d'architecture C4 complète (Diagrammes Mermaidjs Niveaux 1 à 3).
- [x] Établir le comparatif de marché et la validation légale des marques.
- [x] Fixer le nom de code projet : **Project PULSE**.

---

### 1.2 — Modèle de Données & Migrations SQL (⏭️ Prochaine étape)

- [ ] Écrire `migrations/core/000001_create_core_schema.up.sql`
  - [ ] Table `core.sports` (Soccer par défaut, prêt pour multi-sports).
  - [ ] Tables `core.users`, `core.parents_children`, `core.player_profiles`.
  - [ ] Table `core.pools` (Bassins d'âges / Catégories).
  - [ ] Ajouter les instructions `COMMENT ON` pour la documentation automatique.
- [ ] Écrire `migrations/tournament/000001_create_tournament_schema.up.sql`
  - [ ] Tables `tournament.rosters`, `tournament.event_sub_registrations`, `tournament.brackets`.
- [ ] Écrire `migrations/scheduling/000001_create_scheduling_schema.up.sql`
  - [ ] Tables `scheduling.fields`, `scheduling.sub_fields`, `scheduling.matches`, `scheduling.attendances`.
- [ ] Écrire `migrations/finance/000001_create_finance_schema.up.sql`
  - [ ] Tables `finance.event_expenses`, `finance.referee_payments`.
- [ ] Configurer `sqlc.yaml` pour l'isolation stricte des requêtes SQL par module Go.

---

### 1.3 — Socle Backend Go (Modular Monolith)

- [ ] Mettre en place la structure du dépôt Go (`/cmd/monolith`, `/internal/core`, `/internal/tournament`, etc.).
- [ ] Intégrer `golang-migrate` pour l'exécution automatique des migrations SQL par schéma au démarrage.
- [ ] Définir les interfaces Go (`Ports`) pour la communication inter-modules en mémoire.
- [ ] Configurer le serveur HTTP Chi Router et la connexion PostgreSQL via `pgx`.

---

### 1.4 — Socle Frontend React (PULSE UI System)

- [ ] Initialiser le projet Vite + React + TypeScript.
- [ ] Configurer Tailwind CSS, Radix UI et le système de thèmes (Dark/Day mode).
- [ ] Créer les composants de base accessibles WCAG (Button, Modal, Dynamic Data Table, Card).
- [ ] Mettre en place le dictionnaire de termes par sport (Lexique Soccer).

---

## 🔮 PHASES SUIVANTES (Aperçu)

### PHASE 2 : Core Domain & Portail Parent

- [ ] Authentification & Gestion Granulaire des Rôles (RBAC).
- [ ] Espace Parent : Fiche Enfant, Inscription au Bassin U9-U10.
- [ ] Tunnel de Sous-Inscriptions aux tournois optionnels.

### PHASE 3 : Équipes Éphémères & Calendrier Dynamique

- [ ] Outil d'administration : Répartition des jeunes (Draft / Team Builder).
- [ ] Moteur API de consolidation du calendrier familial (Pratiques du bassin + Matchs du roster).
- [ ] Module RSVP Express (Confirmation de présence 1-click).

### PHASE 4 : Logistique Terrains, Brackets & Finances

- [ ] Algorithme d'affectation et de découpage des terrains (11v11 ➔ 2x 9v9 / 4x 7v7).
- [ ] Générateur d'arbres de tournois avec repêchage (Consolidation 3 matchs min. garantis).
- [ ] Dashboard financier par événement (Postes de dépenses vs Sous-inscriptions).

### PHASE 5 : Évaluations & Real-Time

- [ ] Grilles d'évaluation de compétences par calibre (Technique, Tactique, Physique, Mental).
- [ ] Gateway WebSocket Go + Redis Pub/Sub pour la diffusion des scores en direct sur le terrain.

---

## 💡 Reste à Réfléchir / Décisions à Prendre

1. **Paiements :** Choix entre Stripe Connect (ventilation automatique par club) ou Stripe Checkout standard.
2. **Notifications Urgentes :** Priorité aux WebPush ou intégration Twilio SMS pour les annulations de terrain sur le coup (ex: orages).
3. **Identité Finale :** Choix définitif de la marque commerciale (ex: Rostra vs Kipup) avant la mise en production.
4. **Évolutivité :** Préparer le projet pour une future migration vers une architecture microservices si nécessaire (ADR-003).
5. **Multi-Sports :** Définir la stratégie d'ajout de nouveaux sports (soccer, basketball, volleyball, etc.) et l'impact sur le modèle de données et les interfaces utilisateur.
6. **SDLC :** Définir le cycle de vie du développement logiciel, y compris les phases de planification, de conception, de développement, de test et de déploiement.
7. **Sécurité :** Mettre en place des pratiques de sécurité robustes, y compris la gestion des identités et des accès, la protection contre les attaques courantes et la conformité aux réglementations pertinentes (ex: RGPD).
8. Environnement de déploiement : Définir les environnements de développement, de test et de production, ainsi que les stratégies de déploiement continu et d'intégration continue (CI/CD).
# SPÉCIFICATION DU SCHÉMA DE DONNÉES — PROJECT PULSE

## 1. Architecture & Isolation des Schémas SQL

Le système utilise une architecture multi-tenant et multi-sports au sein d'une même instance PostgreSQL, isolée par schémas PostgreSQL.

| Schéma | Domaine | Responsabilité Principale | Statut de Dépendance |
| :--- | :--- | :--- | :--- |
| `core` | Noyau & Entités | Clubs (avec hiérarchie club/association/fédération), utilisateurs, rôles, profils joueurs, sports, positions, bassins d'âge (`pools`), groupes d'entraînement et équipes de saison. | **Kernel Shared** (Référencé par les autres schémas). |
| `tournament` | Compétition | Événements, éligibilité par âge/genre, équipes engagées par un club, alignements (*event rosters*) ponctuels. | Référence `core` et `scheduling` via clés étrangères explicites (FK). |
| `scheduling` | Temporel & Logistique | **Propriétaire unique de toutes les fenêtres temporelles** (saisons, bassins, événements) et du détail fin (terrains, sous-découpage de surfaces, calendriers de matchs/pratiques, présences). | **Kernel Shared** (comme `core`, mais pour le temporel — voir `ADR-008`). |
| `finance` | Comptabilité | Budgets d'événements, frais d'inscription, dépenses, rétribution des arbitres. | Module pair (Pas de FK/JOIN directs avec `tournament`). |
| `evaluation` | Performance | Évaluations détaillées, grilles de compétences, saisie vocale IA (*Voice-to-Eval*). | Module pair (Pas de FK/JOIN directs avec `tournament`). |

## 2. Règles d'Or d'Architecture (ADR-003 & ADR-004)

1. **Relation Core / Modules Pairs (ADR-003 §2.B) :** `core` agit comme noyau partagé. Les schémas satellites (`tournament`, `finance`, `evaluation`) peuvent et doivent utiliser des clés étrangères (FK) directes vers `core` (`core.sports`, `core.pools`, `core.player_profiles`).
2. **`scheduling` : second noyau partagé, pour le temporel (`ADR-008`) :** `scheduling` porte toutes les fenêtres temporelles de la plateforme (saisons, bassins, événements, créneaux). Les modules qui ont besoin d'une date le référencent par FK, plutôt que de dupliquer `start_date`/`end_date` localement.
3. **Isolation entre Modules Pairs (ADR-003) :** Aucune clé étrangère ni `JOIN` SQL direct n'est autorisé entre deux modules métiers pairs (ex: `tournament` vers `finance`). L'intégrité et la communication entre ces domaines s'effectuent exclusivement dans la couche applicative Backend Go via des UUIDs.
4. **Abstraction Multi-Sport (ADR-004) :** Les sports, positions et règles d'alignement sont configurés dynamiquement (`core.sports`, `core.positions`). Aucune modification de schéma SQL n'est requise lors de l'ajout d'un nouveau sport (Soccer, Hockey, Baseball, etc.).

## 3. Structure Détaillée du Schéma `core`

### 3.1 Énumérations

* `core.user_role` : `'SUPER_ADMIN'`, `'CLUB_ADMIN'`, `'TECHNICAL_DIRECTOR'`, `'COACH'`, `'GUARDIAN'`, `'PLAYER'`
* `core.gender_category` : `'MASCULINE'`, `'FEMININE'`, `'MIXED'`
* `core.relationship_type` : `'MOTHER'`, `'FATHER'`, `'LEGAL_GUARDIAN'`, `'OTHER'`
* `core.registration_status` : `'PENDING'`, `'ASSIGNED'`

### 3.2 Tables Principales & Rôles

* **`core.clubs`** : Entités organisationnelles multi-tenants.
* **`core.sports`** : Référentiel global des disciplines avec règles paramétrables (`rules_config` JSONB).
* **`core.users`** : Identités uniques de la plateforme (compte d'accès).
* **`core.user_roles`** : Attributions de rôles par club. Un utilisateur peut cumuler plusieurs rôles (ex: `GUARDIAN` + `COACH`).
  * *Contrainte :* `SUPER_ADMIN` est global à la plateforme (`club_id IS NULL`), les autres rôles sont obligatoirement rattachés à un club. Index d'unicité sur le rôle primaire (`is_primary`).
* **`core.player_profiles`** : Fiches d'identité des athlètes. Lien optionnel vers `core.users` si le joueur a un accès direct.
* **`core.parents_children`** : Table de jonction gérant la tutelle légale et les contacts principaux.

### 3.3 Structuration des Bassins & Positions

* **`core.positions`** : Cartographie des positions spécifiques par sport (`sport_id`, `code`, `name`).
* **`core.player_position_preferences`** : Préférences de postes par joueur, classées par rang.
* **`core.player_ratings`** : Score global d'équilibrage manuel par joueur et par sport (échelle 0 à 100).
* **`core.pools`** : Bassins d'âge et de catégorie (ex: U10F, U12M) rattachés à un club, un sport et une année de saison.
* **`core.pool_divisions`** : Niveaux de compétition au sein d'un bassin (ex: Division 1, Local).
* **`core.pool_registrations`** : Inscription annuelle explicite d'un joueur dans un bassin et une division.

## 4. Structure Détaillée du Schéma `tournament`

### 4.1 Énumérations

* `tournament.roster_type` :
  * `'TRAINING_GROUP'` : Groupe d'entraînement de saison (1:1 avec un bassin `core.pools`).
  * `'SEASON_TEAM'` : Équipe compétitive de saison (peut combiner plusieurs bassins).
  * `'EVENT_ROSTER'` : Alignement éphémère formé spécifiquement pour un événement/tournoi.

### 4.2 Tables & Intégrité

* **`tournament.events`** : Événements/Tournois compétitifs spécifiques rattachés à un sport et une saison.
* **`tournament.rosters`** : Groupes de joueurs.
  * *Contraintes de cohérence :* `EVENT_ROSTER` doit impérativement posséder un `event_id` et des dates nulles (héritées de l'événement). Les types `TRAINING_GROUP` et `SEASON_TEAM` doivent définir leurs propres dates de début/fin sans `event_id`.
* **`tournament.roster_pools`** : Table de jonction reliant un roster aux bassins (`core.pools`) dont il extrait ses joueurs.
* **`tournament.roster_players`** : Appartenance des joueurs aux alignements.
  * *Garantie d'Unicité Événementielle :* Contrainte par index partiel unique `uk_one_roster_per_player_per_event` sur `(event_id, player_id)` interdisant à un joueur d'être inscrit dans plus d'un alignement pour un même événement.

## 5. Matrice des Index & Optimisations Clés

* **Recherche & Filtrage :** Index sur les slugs de clubs (`idx_clubs_slug`), emails utilisateurs, dates de naissance joueurs (`idx_player_profiles_dob`), et codes de sports.
* **Multi-tenant / Performances SQL :** Index composés sur `(club_id, sport_id, season_year)` dans `core.pools` et `(type, season_year)` dans `tournament.rosters`.
* **Contraintes d'Intégrité Performantes :** Les contraintes multi-rôles et l'unicité des joueurs par tournoi sont directement gérées par le moteur PostgreSQL via des index uniques partiels (*Partial Unique Indexes*), éliminant le besoin de triggers SQL lourds.

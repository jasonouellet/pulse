# ADR-008 : Frontière Core / Tournament / Scheduling & Hiérarchie Organisationnelle

## Statut

Accepté — 2026-08-29

## Contexte

La modélisation initiale du schéma `tournament` regroupait trois concepts sous un même type
`roster_type` (`TRAINING_GROUP`, `SEASON_TEAM`, `EVENT_ROSTER`), et les dates (saison, bassin,
événement) étaient dupliquées localement dans chaque schéma. En creusant le flux réel
d'organisation d'un événement multi-clubs, cette frontière s'est révélée incorrecte, et la
structure organisationnelle (`core.clubs`) s'est révélée trop simple pour représenter une
fédération ou une association régionale.

## Décision

### 1. Reclassification `TRAINING_GROUP` / `SEASON_TEAM` → `core`

Ces deux types ne sont pas des concepts de compétition — un groupe d'entraînement et une équipe
de saison existent indépendamment de tout événement/tournoi. Ils rejoignent `core` (aux côtés de
`pools`/`pool_divisions`), comme structure de club persistante à l'année.

`tournament` ne conserve que ce qui est propre à un événement précis : `EVENT_ROSTER` (renommé
alignement — voir §3).

### 2. Flux d'un événement (`tournament`)

Un événement suit un flux en 4 étapes, avec deux inscriptions distinctes qu'il ne faut pas
confondre :

1. **Un club organise un événement** — il définit les modalités (format, règles de classement) et
   l'**éligibilité** : une ou plusieurs plages âge/genre, potentiellement non contiguës sur un
   même événement (ex: U4-U8, U19-F à U12-F, U13-M et plus).
2. **Des clubs s'inscrivent à l'événement** (l'organisateur ou d'autres clubs), en y engageant une
   ou plusieurs **équipes** — une décision de club, par catégorie éligible, référençant le bassin
   source du club engagé. Un événement peut recevoir des équipes engagées par plusieurs clubs
   différents, pas seulement l'organisateur.
3. **Les responsables inscrivent leurs enfants à l'événement** — une sous-inscription séparée de
   l'engagement du club (cohérent avec les "sous-inscriptions à la carte" déjà documentées dans
   `SYSTEM_CONTEXT.md §3`).
4. **Le Directeur Technique construit les alignements** à partir des sportifs inscrits, un par
   équipe engagée, et assigne les entraîneurs pour l'occasion (potentiellement différents du
   coach habituel du groupe d'entraînement).

Le club organisateur gère ensuite l'horaire des matchs, l'assignation des terrains et des
arbitres — voir §4 (`scheduling`).

### 3. Terminologie : "alignement" / "roster"

* **Français (UI, docs) :** *alignement* — remplace "roster" partout où le mot était utilisé en
  français.
* **Anglais / code :** *roster* est conservé — c'est le terme exact pour "liste de joueurs
  rattachés à un groupe", peu importe la durée. *Lineup* est réservé à un concept plus fin, pas
  encore modélisé : qui débute/joue un match précis (titulaires + banc), distinct d'un roster.

### 4. `scheduling` devient propriétaire de toute notion temporelle

`scheduling` ne se limite pas au détail fin (créneaux de matchs/pratiques, terrains) — il devient
l'unique propriétaire de **toutes** les fenêtres temporelles de la plateforme : saisons, fenêtres
de bassin, fenêtres d'événement, et créneaux de matchs/pratiques/terrains. `core.pools` et
`tournament.events` ne stockent plus leurs propres `start_date`/`end_date` en dur ; ils
référencent `scheduling` pour cette information.

**Conséquence sur ADR-003 §2 (isolation entre modules pairs) :** `scheduling` cesse d'être un
module pair strictement isolé de `tournament`/`core` au sens classique — c'est un second type de
dépendance partagée, au même titre que `core`, mais pour le temporel plutôt que pour l'identité.
La règle qui reste absolue : les modules métiers pairs (`tournament`, `finance`, `evaluation`) ne
se référencent **jamais** directement entre eux ; ils peuvent référencer `core` et `scheduling`,
jamais l'un l'autre.

### 5. Rôle "Directeur Technique" = responsable de bassin/alignement

Le rôle qui construit les alignements et assigne les coachs ponctuels est `TECHNICAL_DIRECTOR`
(déjà défini dans `core.user_role`), pas `COACH`. Le Directeur Technique supervise une catégorie
entière à travers le club ; le Coach gère une équipe précise au quotidien. Cette lecture est déjà
cohérente avec `SYSTEM_CONTEXT.md §2`, qui associait déjà "création des événements, répartition
des équipes éphémères" à ce rôle.

### 6. Hiérarchie organisationnelle : club / association / fédération

`core.clubs` gagne :

* `parent_club_id uuid REFERENCES core.clubs (id)` — auto-référencé, nullable, permet une
  imbrication arbitraire (un club peut être rattaché directement à une fédération, ou passer par
  une association intermédiaire).
* `org_type core.organization_type` — enum `CLUB`, `ASSOCIATION`, `FEDERATION`, pour distinguer le
  niveau/rôle de l'organisation sans forcer une profondeur de hiérarchie fixe.

## Conséquences

* `migrations/core/000001_create_core_schema.up.sql` : déplacer `TRAINING_GROUP`/`SEASON_TEAM`
  hors de `tournament.roster_type` vers un nouvel enum/tables `core`, ajouter
  `parent_club_id`/`org_type` sur `core.clubs`.
* `migrations/tournament/000001_create_tournament_schema.up.sql` : `roster_type` ne garde que
  `EVENT_ROSTER` (ou est retiré si un seul type reste — à trancher à l'implémentation), ajoute
  `tournament.event_eligibility` (critères âge/genre, plusieurs lignes possibles par événement) et
  `tournament.team_entries` (club + événement + bassin source).
* `migrations/scheduling/` (nouveau) : porte toutes les dates actuellement en dur dans `core` et
  `tournament`.
* Renommage `roster` → `alignement` dans tout le code/UI en français ; `roster` conservé en
  anglais/code.
* `docs/SCHEMA.md`, `docs/SYSTEM_CONTEXT.md`, `docs/CADRAGE_FONCTIONNEL.md`,
  `docs/ROADMAP_TODO.md` mis à jour en conséquence (voir ces fichiers).

## Alternatives considérées

* **Garder `scheduling` strictement pair, dates dupliquées localement** — rejeté : duplique la
  logique de gestion de dates dans 3+ schémas et risque d'incohérence (ex: la fenêtre d'un
  événement qui ne correspond plus à ses créneaux de matchs).
* **`COACH` construit les alignements** — rejeté : ne correspond pas au flux réel où le Directeur
  Technique supervise plusieurs équipes/coachs à travers une même catégorie lors d'un événement.

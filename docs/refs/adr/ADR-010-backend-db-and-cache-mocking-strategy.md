# [ADR-010] Stratégie de mocking de PostgreSQL et Redis pour le backend Go

* **Statut:** Accepté
* **Décideurs:** Jason Ouellet (Lead Dev), Équipe PULSE
* **Date:** 2026-08-30

## Contexte et énoncé du problème

La couche service et la suite de tests unitaires/d'intégration du backend Go nécessitent l'accès à une base de données PostgreSQL et un cache Redis. Cependant, un environnement Docker n'est pas garanti en permanence sur tous les postes de développement (notamment en environnement restreint WSL2 ou lors de validations rapides sans conteneur).

Comment garantir une suite de tests unitaires ultra-rapide, fiable et exécutable localement en pure Go sans dépendre d'un démon Docker actif ?

## Facteurs décisionnels

* **Autonomie de l'environnement :** Capacité à exécuter `go test ./...` instantanément, sans exiger qu'un démon Docker tourne en arrière-plan.
* **Vitesse d'exécution :** Temps d'initialisation des tests inférieur à quelques millisecondes.
* **Fidélité des comportements :** Simulation précise des commandes Redis (TTL, expiration, structures de clés) et des requêtes SQL/interfaces de driver.
* **Facilité d'intégration :** Compatibilité native avec les packages du projet (`database/sql` / `pgx` et `go-redis`).

## Options envisagées

* **Option 1 (Retenue) :** Mocking pur Go / In-Memory (`pgxmock` / `go-sqlmock` pour Postgres et `alicebob/miniredis` pour Redis)
* **Option 2 :** Conteneurs éphémères réels (`testcontainers-go` avec Docker)
* **Option 3 :** Base de données SQLite en mémoire avec émulation du dialecte Postgres

## Décision retenue

Option choisie : **Option 1**, car elle supprime toute dépendance à l'environnement d'exécution (Docker), garantit un démarrage immédiat des tests (< 10 ms) et permet un développement fluide en pure Go.

### Conséquences positives

* **Zéro dépendance système :** La suite de tests s'exécute sur n'importe quel environnement où le compilateur Go v22+ est installé.
* **Isolation et vitesse :** `alicebob/miniredis` crée un vrai serveur TCP Redis local en mémoire vive à chaque exécution de test sans réutiliser d'état résiduel.
* **Reproductibilité CI/CD :** Les pipelines d'intégration continue peuvent exécuter la suite de tests unitaires sans nécessiter les privilèges Docker (Docker-in-Docker / sockets).

### Conséquences négatives

* **Absence de validation du dialecte SQL réel :** `pgxmock` valide que les requêtes attendues sont appelées et retourne des lignes simulées, mais n'exécute pas le parseur ni les contraintes d'intégrité réelles de PostgreSQL (clés étrangères, déclencheurs, index).
* **Verbosité des mocks SQL :** Il est nécessaire d'écrire explicitement les structures de lignes renvoyées (`sqlmock.NewRows`) pour les requêtes complexes générées par `sqlc`.

## Arguments pour et contre les options

### Option 1 : Pure Go / In-Memory (`pgxmock` + `miniredis`) — *Option retenue*

* **Pour :** Aucun démon Docker requis ; exécution ultra-rapide en mémoire (< 10 ms).
* **Pour :** `miniredis` implémente fidèlement le protocole Redis (TTL, Pub/Sub, Hash, Transactions).
* **Contre :** Les requêtes SQL ne sont pas soumises à un vrai moteur PostgreSQL.

### Option 2 : Conteneurs éphémères (`testcontainers-go`)

* **Pour :** Validation à 100 % de la syntaxe SQL générée par `sqlc` et des migrations.
* **Contre :** Nécessite obligatoirement un démon Docker actif sur l'hôte.
* **Contre :** Temps d'initialisation des conteneurs (~1-3 secondes par suite de tests).

### Option 3 : SQLite en mémoire (`cznic/sqlite`)

* **Pour :** Permet d'exécuter de vraies requêtes SQL en mémoire sans Docker.
* **Contre :** Différences majeures de dialecte avec PostgreSQL (types UUID, JSONB, ENUM et clauses d'UPSERT incompatibles avec `sqlc`).

## Validation de la décision

La décision est validée si:

1. Les tests unitaires de la couche repository/service s'exécutent avec la commande `go test ./...` sans erreur en l'absence de conteneur Docker.
2. `miniredis` intercepte et valide correctement les opérations d'invalidation de cache Redis.
3. Le temps total d'exécution de la suite de tests unitaires reste sous le seuil des 2 secondes.

# ADR 002-C : Choix des Outils Backend Go (sqlc et golang-migrate)

## Statut

Accepté (Précise et amende l'ADR 002-B)

## Contexte

Pour le développement du backend en Go dans **Project PULSE**, nous devions déterminer comment interagir avec la base de données PostgreSQL 16+ tout en respectant l'isolation stricte des schémas SQL (`core`, `tournament`, `scheduling`, etc.).

Deux approches s'offraient à nous :

1. **Un ORM traditionnel (ex: GORM, Ent) :** Abstrait le SQL, mais génère souvent des requêtes inefficaces, masque la complexité des requêtes SQL et facilite accidentellement les jointures inter-schémas interdites par notre architecture.
2. **SQL Pur + Générateur de Code (sqlc) :** Écriture directe des requêtes DDL et SQL nativement, suivie d'une compilation statique vers du code Go typé.

## Décision

Nous adoptons **`sqlc`** pour l'accès aux données et **`golang-migrate`** pour la gestion du schéma SQL.

### 1. Pourquoi `sqlc` (Typesafe SQL Generator) ?

* **Performance Maximale :** Aucun surcoût de réflexion à l'exécution. `sqlc` génère des structures Go natives ultra-rapides au-dessus du driver `pgx/v5`.
* **Sécurité à la Compilation :** `sqlc` valide la syntaxe SQL et les types de colonnes directement contre PostgreSQL au moment de la compilation Go. Si le SQL est invalide, le build échoue.
* **Respect de l'Isolation des Schémas (ADR-003) :** Chaque module Go possède son propre fichier de configuration `sqlc.yaml` et son propre dossier de requêtes SQL. L'outil empêche l'écriture de requêtes croisées non autorisées.
* **Docs-as-Code :** Les commentaires SQL rédigés au-dessus des requêtes sont automatiquement transformés en commentaires de documentation `Godoc`.

### 2. Pourquoi `golang-migrate` ?

* **Migrations Versionnées par Schéma :** Fichiers `.up.sql` et `.down.sql` numérotés pour chaque module (`migrations/core/`, `migrations/tournament/`, etc.).
* **Exécution Automatisée :** Le monolithe Go valide et applique automatiquement les migrations manquantes lors de la séquence de démarrage (`main.go`).

## Conséquences

### Impacts Positifs

* Contrôle total du SQL et optimisation native des index et des colonnes `JSONB`.
* Aucune magie d'ORM : le code Go généré par `sqlc` est simple, lisible et prévisible.
* Intégration parfaite avec l'architecture par Ports & Adapters.

### Contraintes

* Obligation d'écrire le SQL à la main pour chaque requête CRUD (pas de génération automatique de tables à partir de structures Go). Ce choix est assumé pour garantir la qualité de la base de données.

# ADR 003 : Monolithe Modulaire avec Migration Microservices sans Refactoring

## Statut

Accepté

## Contexte

Afin de limiter la complexité opérationnelle initiale, le projet doit démarrer sous forme d'un déploiement unique (Monolithe). Toutefois, pour absorber les pics de charge (ex: calculs lourds d'arbres de tournois) et permettre l'isolation future, l'architecture doit permettre d'extraire n'importe quel module vers un microservice conteneurisé indépendant sans modifier la logique métier.

## Décisions

1. **Modèle de Code : Monolithe Modulaire (Modular Monolith)**
   * Le code est structuré dans un Monorepo Go par sous-domaines stricts (`internal/core`, `internal/tournament`, etc.).
   * Aucune dépendance directe entre les paquets métiers. Toute communication inter-modules passe par des **Ports (Interfaces Go)**.
2. **Isolation des Données (Database-per-module logique)**
   * Utilisation d'un seul serveur PostgreSQL au départ, mais avec des **Schémas SQL séparés par module** (`core`, `tournament`, `scheduling`, `finance`, `evaluation`).
   * Interdiction absolue de requêtes JOIN entre les tables de modules différents.
2.B. **Modèle de dépendance à deux niveaux (Shared Kernel)**
   * `core` est un **noyau partagé** : il porte les entités transversales (`sports`, `users`, `player_profiles`, `pools`, `positions`, etc.) dont **tous** les autres modules dépendent légitimement. Les autres schémas (`tournament`, `scheduling`, `finance`, `evaluation`) peuvent référencer `core` par **clé étrangère SQL réelle** (`REFERENCES core.xxx (id)`), pas seulement par UUID validé en application — l'intégrité référentielle en base est préférable tant que le coût de migration vers un microservice n'est pas encore engagé.
   * Ce que l'isolation interdit reste **les dépendances entre modules pairs** : `tournament` ne doit jamais référencer directement `finance` ou `scheduling`, et vice-versa. Toute communication entre modules pairs passe par les Ports (Go interfaces), jamais par une FK ou un JOIN cross-schéma entre eux.
   * Conséquence assumée pour l'extraction microservice future : si `tournament` est extrait avant `core`, les FK vers `core` devront être remplacées par de la validation applicative (appel au service `core` ou événements) au moment de l'extraction — c'est un coût connu et accepté, pas une raison d'éviter les FK aujourd'hui.
3. **Inversion de Dépendance & Injection dynamique (Adapters)**
   * Chaque module expose un adaptateur `In-Memory` (pour le mode Monolithe) et un adaptateur `gRPC/Event` (pour le mode Microservice).
   * L'instanciation de l'adaptateur se fait au démarrage du conteneur via la configuration/variables d'environnement.
4. **Conteneurisation (Docker)**
   * Utilisation de _Multi-stage Builds_ Docker permettant de produire soit le binaire global (`cmd/backend`), soit un microservice spécifique (`cmd/services/X`).

## Conséquences

* Zéro surcoût réseau au départ, déploiement initial ultra-simple, flexibilité totale pour distribuer les charges plus tard.
* Discipline d'équipe obligatoire (interdiction des JOINs SQL inter-schémas).

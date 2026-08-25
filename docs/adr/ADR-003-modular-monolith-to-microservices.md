# ADR 003 : Monolithe Modulaire avec Migration Microservices sans Refactoring

## Statut
Accepté

## Contexte
Afin de limiter la complexité opérationnelle initiale, le projet doit démarrer sous forme d'un déploiement unique (Monolithe). Toutefois, pour absorber les pics de charge (ex: calculs lourds d'arbres de tournois) et permettre l'isolation future, l'architecture doit permettre d'extraire n'importe quel module vers un microservice conteneurisé indépendant sans modifier la logique métier.

## Décisions
1. **Modèle de Code : Monolithe Modulaire (Modular Monolith)**
   - Le code est structuré dans un Monorepo Go par sous-domaines stricts (`internal/core`, `internal/tournament`, etc.).
   - Aucune dépendance directe entre les paquets métiers. Toute communication inter-modules passe par des **Ports (Interfaces Go)**.
2. **Isolation des Données (Database-per-module logique)**
   - Utilisation d'un seul serveur PostgreSQL au départ, mais avec des **Schémas SQL séparés par module** (`core`, `tournament`, `scheduling`, `finance`, `evaluation`).
   - Interdiction absolue de requêtes JOIN entre les tables de modules différents.
3. **Inversion de Dépendance & Injection dynamique (Adapters)**
   - Chaque module expose un adaptateur `In-Memory` (pour le mode Monolithe) et un adaptateur `gRPC/Event` (pour le mode Microservice).
   - L'instanciation de l'adaptateur se fait au démarrage du conteneur via la configuration/variables d'environnement.
4. **Conteneurisation (Docker)**
   - Utilisation de *Multi-stage Builds* Docker permettant de produire soit le binaire global (`cmd/monolith`), soit un microservice spécifique (`cmd/services/X`).

## Conséquences
- Zéro surcoût réseau au départ, déploiement initial ultra-simple, flexibilité totale pour distribuer les charges plus tard.
- Discipline d'équipe obligatoire (interdiction des JOINs SQL inter-schémas).

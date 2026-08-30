# ARCHITECTURE TECHNIQUE — PROJECT PULSE.OS

## 1. Contexte Système (Niveau 1)

L'application Soccer Manager ERP permet aux clubs, entraîneurs et parents de gérer la logistique sportive (bassins, rosters éphémères, tournois, terrains, finances et évaluations).
Le système communique avec des services tiers pour le paiement (Stripe), les notifications (SMS/Push) et la géolocalisation des terrains.

Ce diagramme illustre les interactions entre les utilisateurs, le système central et les services externes.

```mermaid
C4Context
    title Diagramme de Contexte Système (Niveau 1) - Soccer Manager ERP

    Person(admin, "Administrateur Club", "Gère la structure, les terrains, les événements et le budget.")
    Person(coach, "Entraîneur / Staff", "Gère les présences, les alignements et les évaluations.")
    Person(parent, "Parent / Tuteur", "Gère les inscriptions, les tournois à la carte et l'horaire familial.")
    Person(player, "Sportif (Joueur)", "Consulte ses matchs, ses statistiques et ses objectifs.")

    System(soccer_erp, "Soccer Manager ERP", "Plateforme SaaS de gestion de soccer multi-niveaux (bassins, rosters éphémères, tournois, terrains, finances).")

    System_Ext(stripe, "Système de Paiement", "Stripe - Gestion des cotisations et sous-inscriptions.")
    System_Ext(notifications, "Service de Notification", "Twilio / WebPush - Envoi d'alertes, SMS et rappels.")
    System_Ext(maps, "Service de Cartographie", "Google / Apple Maps - Géolocalisation des terrains et complexes.")

    Rel(admin, soccer_erp, "Gère et configure", "HTTPS / Web UI")
    Rel(coach, soccer_erp, "Saisit présences et évaluations", "HTTPS / PWA")
    Rel(parent, soccer_erp, "Consulte horaire et inscrit enfants", "HTTPS / Mobile UI")
    Rel(player, soccer_erp, "Consulte son profil et ses matchs", "HTTPS / Mobile UI")

    Rel(soccer_erp, stripe, "Traite les paiements", "REST API")
    Rel(soccer_erp, notifications, "Déclenche les alertes", "REST API")
    Rel(soccer_erp, maps, "Affiche les itinéraires", "REST API")
```

## 2. Conteneurs (Niveau 2)

* **Frontend SPA :** React.js + TypeScript (Mobile-First, PWA, Dark/Day Mode, WCAG 2.1 AA).
* **Backend API :** Go (Golang) sous forme de Monolithe Modulaire.
* **Base de Données :** PostgreSQL 16+ avec un schéma dédié par module.
* **Cache & Direct :** Redis (Pub/Sub WebSockets pour les scores et cache d'horaires).

```mermaid
C4Container
    title Diagramme de Conteneurs (Niveau 2) - Sport Manager ERP

    Person(user, "Utilisateur (Tous rôles)", "Accède à la plateforme via le navigateur ou le mobile.")

    Container_Boundary(bnd_erp, "Sport Manager ERP System") {
        Container(spa, "Single Page Application (SPA)", "React.js, TypeScript, Tailwind", "Interface web/mobile universelle, accessible WCAG 2.1 AA, support Day/Dark mode.")
        Container(api, "Backend API Server", "Go (Golang)", "Monolithe modulaire traitant la logique métier, les brackets et la planification.")
        ContainerDb(postgres, "Base de Données Principale", "PostgreSQL 16+", "Stockage relationnel isolé par schémas SQL (core, tournament, scheduling, finance, evaluation, etc.).")
        ContainerDb(redis, "Cache & Broker Temps Réel", "Redis 7+", "Gestion des sessions, verrous distribués sur terrains, cache HTTP et Pub/Sub WebSockets.")
    }

    System_Ext(stripe, "Stripe API", "Paiements")
    System_Ext(push, "Push Notification Service", "SMS / WebPush")

    Rel(user, spa, "Utilise", "HTTPS / PWA")
    Rel(spa, api, "Effectue des requêtes API", "JSON / HTTPS")
    Rel(spa, api, "Reçoit les scores en direct", "WebSockets")

    Rel(api, postgres, "Lit / Écrit", "pgx / SQL Native")
    Rel(api, redis, "Met en cache & Pub/Sub", "Redis Protocol")

    Rel(api, stripe, "Initialise paiements", "HTTPS / REST")
    Rel(api, push, "Envoie alertes", "HTTPS / REST")
```

## 3. Composants Backend (Niveau 3)

Chaque module dans `/internal/` respecte l'Architecture Hexagonale:

* **`domain` :** Modèles et règles métier purement Go (aucune dépendance).
* **`ports` :** Interfaces de communication d'entrée et de sortie.
* **`adapters` :** Implémentations concrètes (sqlc/PostgreSQL, Redis, In-Memory ou gRPC).

```mermaid
C4Component
    title Diagramme de Composants Backend (Niveau 3) - Monolithe Modulaire Go

    Container_Boundary(bnd_go, "Backend Go (Modular Monolith)") {
        Component(router, "HTTP Router & WS Gateway", "Chi Router", "Achemine les requêtes REST et gère les connexions WebSockets.")

        Component(mod_core, "Module Core & Users", "Go Package (/internal/core)", "Gestion des utilisateurs, liens parent-enfant, RBAC, bassins d'âges.")
        Component(mod_tourn, "Module Tournament", "Go Package (/internal/tournament)", "Gestion des brackets, rosters éphémères, règle des 3 matchs min.")
        Component(mod_sched, "Module Scheduling", "Go Package (/internal/scheduling)", "Terrains, sous-découpage (11v11->7v7), calendriers, RSVP.")
        Component(mod_fin, "Module Finance", "Go Package (/internal/finance)", "Budgets d'événements, sous-inscriptions, paies arbitres.")
        Component(mod_eval, "Module Evaluation", "Go Package (/internal/evaluation)", "Grilles d'évaluations par calibre, fiches de suivi.")

        Component(ports, "Interfaces Inter-Modules (Ports)", "Interfaces Go", "Assure l'isolation stricte et la communication in-memory entre modules.")
    }

    ContainerDb(postgres, "PostgreSQL 16+", "Schémas SQL Cloisonnés")

    Rel(router, mod_core, "Appelle")
    Rel(router, mod_tourn, "Appelle")
    Rel(router, mod_sched, "Appelle")
    Rel(router, mod_fin, "Appelle")
    Rel(router, mod_eval, "Appelle")

    Rel(mod_core, ports, "Utilise")
    Rel(mod_tourn, ports, "Utilise")
    Rel(mod_sched, ports, "Utilise")
    Rel(mod_fin, ports, "Utilise")
    Rel(mod_eval, ports, "Utilise")

    Rel(mod_core, postgres, "Requêtes SQL", "schema 'core'")
    Rel(mod_tourn, postgres, "Requêtes SQL", "schema 'tournament'")
    Rel(mod_sched, postgres, "Requêtes SQL", "schema 'scheduling'")
    Rel(mod_fin, postgres, "Requêtes SQL", "schema 'finance'")
    Rel(mod_eval, postgres, "Requêtes SQL", "schema 'evaluation'")
```

## 4. Composants Frontend SPA (Niveau 3 — Architecture React)

Structure interne de l'application Single Page Application (SPA) React en architecture modulaire par fonctionnalités (Feature-Driven Folder Structure).

```mermaid
C4Component
    title Diagramme de Composants Frontend SPA (Niveau 3)

    Container_Boundary(bnd_react, "Frontend Single Page Application (React + Vite)") {
        Component(router_ui, "React Router & Layout Shell", "React Router v6", "Gère la navigation, le mode Day/Dark et le thème multi-sport actif (Soccer, Hockey, Baseball).")
        Component(auth_store, "Auth & Session Store", "Zustand Store", "Gère l'état d'authentification, le JWT et le rôle actif (UI Role Switcher).")
        Component(theme_context, "Sport Theme Provider", "React Context", "Injecte dynamiquement les variables CSS de sport (`data-sport` & `data-mode`).")
        Component(feat_roster, "Roster & Pool Feature", "React Components / TanStack Query", "Interface Drag & Drop d'équilibrage d'alignements et gestion des bassins d'âge.")
        Component(feat_eval, "Voice-to-Eval Feature", "React Web Audio API", "Capture audio du terrain, pré-visualisation et envoi de la dictée vocale.")
        Component(api_client, "Huma v2 API Client", "Fetch API / TanStack Query", "Couche client HTTP typée communiquant avec l'API REST backend Go.")
    }
    Container(api_server, "Backend API Server", "Go (Golang)", "Huma v2 REST Server")

    Rel(router_ui, auth_store, "Lit le rôle actif")
    Rel(router_ui, theme_context, "Applique le style du sport")
    Rel(router_ui, feat_roster, "Affiche la route /rosters")
    Rel(router_ui, feat_eval, "Affiche la route /evaluations")

    Rel(feat_roster, api_client, "Requêtes d'alignement", "JSON / HTTPS")
    Rel(feat_eval, api_client, "Envoie fichiers audio & notes", "Multipart / HTTPS")

    Rel(api_client, api_server, "Appels API REST / WebSockets", "HTTPS / WSS")
```

## 3. Composants de Stockage & Persistance (Niveau 3 — PostgreSQL & Redis)

### A. Composants PostgreSQL (Schémas Isolés & Moteur de Données)

```mermaid
C4Component
    title Diagramme de Composants - Base de Données PostgreSQL (Niveau 3)

    Container_Boundary(bnd_pg, "Instance PostgreSQL 16+") {
        Component(sch_core, "Schéma 'core'", "PostgreSQL Schema (Kernel)", "Clubs, utilisateurs, rôles, profils joueurs, sports, positions et pools.")
        Component(sch_tourn, "Schéma 'tournament'", "PostgreSQL Schema", "Événements, rosters, équipes de saison et règles d'unicité (uk_one_roster_per_player_per_event).")
        Component(sch_sched, "Schéma 'scheduling'", "PostgreSQL Schema", "Terrains, sous-découpage de surfaces, calendriers, présences et RSVP.")
        Component(sch_fin, "Schéma 'finance'", "PostgreSQL Schema", "Budgets d'événements, cotisations, sous-inscriptions et paies arbitres.")
        Component(sch_eval, "Schéma 'evaluation'", "PostgreSQL Schema", "Fiches d'évaluations athlétiques, grilles par calibre et historique des dictées vocales.")
    }

    Container(api_go, "Backend API Server", "Go (Golang)", "Gère l'accès aux schémas via pgxpool / sqlc.")

    Rel(api_go, sch_core, "FKs directes depuis les autres schémas", "SQL Native / pgx")
    Rel(api_go, sch_tourn, "FKs vers 'core' uniquement (ADR-003)", "SQL Native / pgx")
    Rel(api_go, sch_sched, "Module isolé (Pas de FKs inter-pairs)", "SQL Native / pgx")
    Rel(api_go, sch_fin, "Module isolé (Pas de FKs inter-pairs)", "SQL Native / pgx")
    Rel(api_go, sch_eval, "Module isolé (Pas de FKs inter-pairs)", "SQL Native / pgx")
```

### B. Composants Redis (Cache, Sessions & Temps Réel)

```mermaid
C4Component
    title Diagramme de Composants - Redis 7+ (Niveau 3)

    Container_Boundary(bnd_redis, "Instance Redis 7+") {
        Component(red_session, "Session & Auth Store", "Redis Key-Value (TTL)", "Stockage temporaire des tokens JWT révoqués, refresh tokens et sessions d'utilisateurs.")
        Component(red_lock, "Distributed Lock Manager", "Redlock / Redis Keys", "Verrous distribués pour la réservation de terrains et la modification simultanée des rosters.")
        Component(red_cache, "API Query Cache", "Redis Hash / Strings", "Mise en cache des données de référence peu changeantes (sports, positions, calendriers de matchs).")
        Component(red_pubsub, "WebSockets Pub/Sub Broker", "Redis Pub/Sub Channels", "Acheminement en temps réel des scores de matchs, changements d'alignement et notifications terrain.")
    }

    Container(api_go, "Backend API Server", "Go (Golang)", "Client Rueidis / Go-Redis")
    Container(spa_react, "Frontend SPA", "React / WebSocket", "Reçoit les événements temps réel")

    Rel(api_go, red_session, "Lit / Valide sessions", "RESP Protocol")
    Rel(api_go, red_lock, "Acquiert / Libère verrous", "RESP Protocol")
    Rel(api_go, red_cache, "Invalide & lit le cache", "RESP Protocol")
    Rel(api_go, red_pubsub, "Publie événements (Match Events)", "RESP Protocol")

    Rel(red_pubsub, spa_react, "Diffuse via Gateway WS", "WebSockets (WSS)")
```

## [TODO ARCHITECTURE C4]

* [ ] Spécifier les contrats protobuf/gRPC si un module doit être extrait en microservice.
* [ ] Modéliser le schéma de données détaillé du module `evaluation` (grilles dynamiques JSONB).

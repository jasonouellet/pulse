# Rules of Architecture & Package Structure — Project PULSE

Ce document définit les règles de conception obligatoires pour le backend Go du projet PULSE. Ces règles garantissent le respect strict du patron **Architecture Hexagonale (Ports & Adaptateurs)**, la découplabilité des modules et la testabilité sans dépendance d'infrastructures externes.

## **Règle 1 : Règle de Dépendance Unidirectionnelle (Clean Architecture)**

Le sens des dépendances doit impérativement pointer vers l'intérieur :

$$\text{Adapters} \longrightarrow \text{Ports} \longleftarrow \text{Domain}$$

* **Ports (`internal/<module>/ports/`) :** Contiennent uniquement les interfaces Go pures et les structures de transfert (DTOs). Ils ne doivent **jamais** importer de packages d'infrastructure (ex: `github.com/jackc/pgx`, `github.com/danielgtaylor/huma`, `database/sql`).
* **Adaptateurs (`internal/<module>/adapters/`) :** Importent et implémentent les interfaces définies dans `ports/`.
* Aucune dépendance cyclique ne sera tolérée entre les packages.

## **Règle 2 : Emplacement Unique et Canonique de la Couche HTTP**

Afin d'éviter la fragmentation du code HTTP et la création de répertoires parallèles incohérents :

* **Emplacement exclusif :** Tous les handlers REST/OpenAPI, les types de requête/réponse HTTP (transport DTOs) et les middlewares HTTP doivent obligatoirement résider sous :  
  `internal/<module>/adapters/http/`[cite: 1]
* **Interdiction :** Il est strictement interdit de créer des dossiers racines ou secondaires comme `internal/api/`, `internal/routes/` ou `internal/controllers/`[cite: 1].

## **Règle 3 : Découplage Systématique par Interfaces**

* **Injection de dépendances :** Tout service métier ou adaptateur HTTP doit consommer ses dépendances de persistance ou de services externes à travers des interfaces[cite: 1].
* **Constructeurs :** Les constructeurs d'adaptateurs doivent toujours prendre en paramètre des types d'interfaces (ex: `func NewTeamHandler(repo ports.TeamRepository)`) et jamais des implémentations concrètes[cite: 1].

## **Règle 4 : Segregation du Code d'Infrastructure et de Plateforme**

* **`pkg/` (Partagé & Public) :** Contient le code technique réutilisable par tous les modules sans logique métier (ex: pool PostgreSQL `pkg/database`[cite: 1], télémétrie `pkg/observability`[cite: 1]).
* **`internal/platform/` (Infrastructures transversales) :** Contient les adaptateurs d'exécution et de plateforme spécifiques à l'application (ex: client Redis `internal/platform/redis/`[cite: 1], gestionnaire de persistance des sessions `internal/session/`[cite: 1]).

## **Règle 5 : Arborescence Modulaire Standardisée**

Chaque nouveau module fonctionnel (ex: `core`, `scheduling`, `tournament`)[cite: 1] doit suivre rigoureusement l'arborescence ci-dessous :

```text
internal/
├── <module>/                    # Nom du module fonctionnel (ex: core, tournament)
│   ├── ports/                   # Contrats d'entrée/sortie (Interfaces, DTOs métiers)
│   └── adapters/                # Implémentations techniques des ports
│       ├── http/                # Handlers Huma/Chi, middlewares et DTOs de transport
│       └── postgres/            # Requêtes SQL et repositories pgx/sqlc
├── platform/                    # Services d'infrastructure de la plateforme
│   ├── redis/                   # Client Redis (go-redis / miniredis)
│   └── session/                 # Store et gestion de sessions utilisateur
└── testutil/                    # Helpers, fixtures et mocks réutilisables pour les tests
```

## Règle 6 : Stratégie de Testabilité Isolée

* Tous les tests unitaires et d'intégration de couche doivent pouvoir s'exécuter localement via go test ./... sans nécessiter de démon Docker actif.
* Mocks autorisés:
  * PostgreSQL : Utilisation de pgxmock pour simuler le comportement du pool SQL
  * Redis : Utilisation de alicebob/miniredis pour simuler l'instance in-memory en pure Go.

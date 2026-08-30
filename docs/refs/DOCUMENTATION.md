# STRATÉGIE DOCUMENTAIRE — PROJECT PULSE

## 1. Principes Directeurs (Docs-as-Code)

La documentation de **Project PULSE** est traitée comme du code source:

* Elle réside directement dans le dépôt Git (sous `/docs`).
* Elle est écrite exclusivement en **Markdown** et validée en CI via `markdownlint-cli2`.
* Sa génération est automatisée à partir du code (Go, SQL, OpenAPI) pour éviter toute dette documentaire.

## 2. Documentation Technique (Pour les Développeurs)

L'objectif est d'automatiser 90 % de la documentation technique pour qu'elle reste toujours synchronisée avec la réalité du code.

### A. Spécification d'API (Backend Go)

* **Outil :** OpenAPI 3.1 généré nativement par **Huma v2**.
* **Accès local :** `http://localhost:8888/docs` (Swagger UI / RapiDoc).
* **Exportation CI :** Un script d'intégration continue exporte `docs/refs/openapi.json` à chaque build pour fournir le contrat d'API aux développeurs frontend.

### B. Base de Données & Schémas SQL

* **Outil :** **tbls** (ou **schemaspy**).
* **Fonctionnement :** Analyse l'instance PostgreSQL à partir des commentaires SQL (`COMMENT ON TABLE` et `COMMENT ON COLUMN`) pour générer automatiquement les diagrammes ERD (SVG) et la description Markdown de tous les schémas isolés (`core`, `tournament`, `scheduling`, ...).

### C. Architecture & Décisions Structurantes

* **Diagrammes C4 :** Documentés et maintenus dans `ARCHITECTURE_C4.md` (Niveaux 1 à 3 : Contexte, Conteneurs, Composants Backend/Persistance).
* **Registres de Décision (ADR) :** Chaque choix technique majeur fait l'objet d'un fichier Markdown sous `/docs/refs/adr/` suivant le format **MADR 4.0.0**.

### D. Design System & Frontend (React)

* **Outil :** **Storybook** / **TypeDoc**.
* **Contenu :** Documentation des composants React, du sélecteur de rôle (`UI Role Switcher`), et du système de thèmes dynamiques (`data-sport` et `data-mode`).

## 3. Documentation Utilisateur (Framework Diátaxis)

Pour les clubs, administrateurs, entraîneurs et parents, la documentation est rédigée selon la méthodologie [Diátaxis](https://diataxis.fr/).
Elle sépare strictement les contenus en 4 quadrants selon l'objectif de l'utilisateur:

| Quadrant                | Orientation                          | Objectif                                                                       | Type de Contenu                                                                                                                  |
| :---------------------- | :----------------------------------- | :----------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| **1. Tutoriels**        | Apprentissage *(Pratique)*           | Guider les débutants pas à pas vers un premier résultat concret.               | Créer un premier club, configurer un bassin U10F, réaliser une première dictée vocale.                                           |
| **2. Guides Pratiques** | Résolution de Problèmes *(Pratique)* | Fournir une recette directe pour accomplir une tâche métier précise.           | Découper un terrain 11v11 en deux surfaces 7v7, valider un consentement parental (< 14 ans / Loi 25), exporter une fiche joueur. |
| **3. Référence**        | Information *(Théorique)*            | Décrire de façon factuelle, neutre et exhaustive le fonctionnement du système. | Dictionnaire des permissions RBAC, grilles des calibres par sport, mots-clés reconnus par l'IA vocale.                           |
| **4. Explications**     | Compréhension *(Théorique)*          | Expliquer le contexte et les raisons des choix métier ou légaux.               | Pourquoi la règle d'unicité par événement ?, comment les PII sont masquées avant transmission aux modèles IA.                    |

## 4. Automatisation & Pipeline CI/CD

Les processus d'actualisation de la documentation sont intégrés au flux GitHub Actions / pre-commit :

1. **Pre-commit :** Validation de la syntaxe des fichiers Markdown via `markdownlint-cli2`.
2. **Commit / Build Backend :** Exécution automatique de la génération du schéma OpenAPI 3.1.
3. **Migration SQL :** Régénération automatique des diagrammes de la base de données via `tbls`.
4. **Release :** Compilation du site de doc utilisateur via un générateur de site statique (**Starlight / Astro / Mkdocs**) hébergé sur sous-domaine dédié (`docs.pulse.app`).

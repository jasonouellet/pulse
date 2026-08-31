# [ADR-0004] Standardisation de l'Architecture Hexagonale et de la Structure des Packages Backend Go

* **Statut:** Accepté
* **Décideurs:** Jason Ouellet (Lead Dev), Équipe PULSE
* **Date:** 2026-08-30

## Contexte et énoncé du problème

Au fur et à mesure de la croissance du backend Go pour Project PULSE, le risque de dérive architecturale s'est accru. Pour maintenir la séparation des responsabilités et éviter les cycles de dépendances, nous devons formaliser un patron d'architecture canonique et fixer des règles strictes sur la structure des packages.

## Facteurs décisionnels

* **Isolation métier :** Le cœur du domaine doit rester indépendant des frameworks Web (Chi, Huma), des ORMs et des détails d'infrastructure.
* **Testabilité :** Les composants métier et la persistance doivent pouvoir être testés de manière isolée sans conteneurs externes ni serveur HTTP.
* **Cohérence modulaire :** Chaque module fonctionnel doit suivre la même arborescence prévisible.
* **Prévention du couplage :** Les dépendances doivent impérativement pointer vers l'intérieur (Règle de dépendance de la Clean Architecture).

## Options envisagées

* **Option 1 (Retenue) :** Architecture Hexagonale (Ports & Adaptateurs) avec layout par module sous `internal/`.
* **Option 2 :** Layout traditionnel par couche globale (`controllers/`, `services/`, `models/`, `repositories/`).
* **Option 3 :** Layout piloté par les fonctionnalités sans séparation claire entre ports et adaptateurs (*Feature Folders*).

## Décision retenue

Option choisie : **Option 1**, car elle garantit l'étanchéité du domaine métier, prévient les cycles d'importation Go et standardise l'injection de dépendances via des interfaces pures (`ports/`).

### Conséquences positives

* **Invariabilité du Domaine :** Remplacer un framework Web (ex: passer de Chi à Fiber) ou un moteur de stockage (ex: passer de PostgreSQL à MongoDB) n'impacte aucunement la logique métier ni les ports.
* **Découpage strict :** Tous les handlers et middlewares HTTP sont regroupés sous les adaptateurs HTTP (`adapters/http/`), éliminant la confusion des dossiers d'entrée multiples[cite: 1].
* **Facilité de refactorisation :** L'usage d'interfaces explicites (`ports/`) facilite la simulation via des mocks (`pgxmock`[cite: 1], `miniredis`[cite: 1]).

### Conséquences négatives

* **Légère verbosité initiale :** Nécessite de définir explicitement des interfaces Go (ports)[cite: 1] et des structures DTO de transfert[cite: 1] au lieu de réutiliser directement les structs de base de données.

## Validation de la décision

La décision est validée si :
1. Aucun fichier sous `internal/<module>/ports/`[cite: 1] n'importe de dépendance HTTP (Huma, Chi)[cite: 1] ou SQL (`pgx`, `sql`)[cite: 1].
2. Tous les handlers et middlewares HTTP sont situés sous `internal/<module>/adapters/http/`[cite: 1].
3. La suite de tests s'exécute de manière isolée via la commande `go test ./...`[cite: 1].

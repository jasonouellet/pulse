# SPÉCIFICATION DU SCHÉMA DE DONNÉES — PROJECT PULSE

## Isolation des Schémas SQL

| Schéma       | Domaine                  | Responsabilité principale                                        |
| :----------- | :----------------------- | :--------------------------------------------------------------- |
| `core`       | Utilisateurs & Structure | Gestion des comptes, liens familiaux, bassins (`pools`), sports. |
| `tournament` | Compétition              | Brackets, sous-inscriptions, rosters éphémères.                  |
| `scheduling` | Logistique               | Terrains, sous-découpage (11v11 -> 7v7), calendriers, présences. |
| `finance`    | Comptabilité             | Budgets d'événements, dépenses, rétribution arbitres.            |

## Règle d'Or d'Architecture

- **Aucune clé étrangère (FK) inter-schémas.**
- L'intégrité entre les schémas `core.pools` et `tournament.rosters` est maintenue applicativement dans le backend Go via des UUIDs.

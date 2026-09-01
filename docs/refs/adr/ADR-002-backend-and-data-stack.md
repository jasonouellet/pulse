# [ADR-002]: Choix de la Pile Backend Compilée et Stratégie de Données Hybride

## Statut

Accepté

## Contexte

Recherche de performances élevées, d'une faible empreinte mémoire et d'une sécurité à la compilation pour les moteurs de calculs (calendriers, brackets de tournois), tout en garantissant un support hautement concurrent pour le rafraîchissement des scores sur le terrain.

## Décisions

1. **Framework Backend :** Go (Golang) avec architecture hexagonale (Clean Architecture).
   * Utilisation de `Chi` comme routeur HTTP/REST.
   * Utilisation de WebSockets pour les mises à jour de matchs en temps réel.
2. **Base de Données Principale :** PostgreSQL 16+ (Accès typé via `sqlc`).
   * Maintien des garanties ACID pour les finances et inscriptions.
   * Utilisation de colonnes `JSONB` pour les structures de données flexibles (ex: grilles d'évaluations).
3. **Cache & File d'Attente :** Redis 7+
   * Cache de premier niveau pour les calendriers, classements de tournois et fiches de présence.
   * Gestion du Pub/Sub pour la diffusion en temps réel des scores.

## Conséquences

* Temps de réponse de l'API extrêmement rapides.
* Consommation mémoire du serveur backend réduite (< 50 Mo de RAM au repos).
* Nécessité de maintenir un DTO structuré entre le backend Go et le frontend React.

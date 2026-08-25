# ADR 004 : Modèle d'Abstraction Multi-Sports avec Focus Initial Soccer

## Statut

Accepté

## Contexte

L'application vise à terme la gestion de multiples sports collectifs (hockey, basket-ball, volleyball, etc.). Toutefois, la Phase 1 doit être livrée rapidement avec une spécialisation complète sur le soccer. Nous devons éviter de sur-complexifier le code tout en empêchant le couplage dur au vocabulaire ou aux règles exclusives du soccer.

## Décisions

1. **Entité "Sport" comme Clé Pivot :** Toutes les structures de compétition, catégories, terrains et grilles d'évaluation incluent une référence `sport_id` (définie dans `core.sports`).
2. **Phase 1 (Soccer-Centric) :** La valeur par défaut `SOCCER` est injectée automatiquement dans l'application. Les interfaces utilisateur afficheront les termes spécifiques au soccer (ex: "Match", "Terrain", "Mi-temps").
3. **Configurations via JSONB :** Les spécificités d'un sport (formats de jeu, durée des périodes, types de scores/points) sont stockées sous forme de paramètres flexibles JSONB pour éviter de modifier la structure de la base de données lors de l'ajout d'un nouveau sport en Phase 2.

## Conséquences

- Zéro refactoring de la base de données lorsque le 2e sport sera ajouté.
- Nécessité de maintenir un dictionnaire de termes/lexique par sport pour la couche d'affichage React.

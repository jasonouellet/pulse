# ANALYSE COMPARATIVE DES FONCTIONNALITÉS — PROJECT PULSE

## 1. Positionnement Stratégique sur le Marché

Le marché actuel des logiciels de gestion sportive est scindé en deux catégories :

1. **Les outils d'associations/ligues traditionnelles (ex: Spordle/Retro) :** Lourds, rigides, conçus pour la gestion administrative et les formulaires d'affiliation, mais très peu adaptés à la logistique terrain quotidienne.
2. **Les applications d'équipe (ex: TeamSnap, SportEasy) :** Conçues pour une seule équipe fixe à la fois. Dès qu'un club veut gérer un bassin de joueurs, des rotations de tournois ou du multi-niveaux, ces outils atteignent leurs limites.

**Project PULSE** comble ce vide avec le modèle **"Pools & Dynamic Rosters"** (Bassins & Alignements Éphémères).

---

## 2. Matrice Comparative des Fonctionnalités

| Domaine Fonctionnel         | Fonctionnalité Clé                                      | Spordle / Retro |  TeamSnap  |     SportEasy     |  **Project PULSE**   |
| :-------------------------- | :------------------------------------------------------ | :-------------: | :--------: | :---------------: | :------------------: |
| **Gestion des Alignements** | Équipes fixes par saison                                |       ✅        |     ✅     |        ✅         |          ✅          |
|                             | **Bassins d'âges (Pools) & Rosters Éphémères**          |       ❌        |     ❌     |        ❌         |   **✅ Exclusif**    |
|                             | **Outil de Répartition (Draft / Team Builder)**         |       ❌        |     ❌     |        ❌         |     **✅ Natif**     |
| **Inscriptions & Finances** | Inscription saisonnière globale                         |       ✅        |     ✅     |        ⚠️         |          ✅          |
|                             | **Sous-inscriptions à la carte (Tournois #1, #3)**      |       ❌        |     ❌     |        ❌         |     **✅ Natif**     |
|                             | **Budget & Rémunérations par événement**                |       ❌        |     ❌     |        ❌         |     **✅ Natif**     |
| **Logistique & Terrains**   | Horaires et réservations simples                        |       ✅        |     ✅     |        ✅         |          ✅          |
|                             | **Sous-découpage dynamique des terrains (11v11 ➔ 7v7)** |       ❌        |     ❌     |        ❌         |     **✅ Natif**     |
|                             | **Horaire familial dynamique consolidé**                |   ⚠️ Partiel    | ⚠️ Partiel |    ⚠️ Partiel     |   **✅ Optimisé**    |
| **Compétition & Tournois**  | Formats standards (À la ronde / Élimination)            |       ✅        | ⚠️ Basique |        ❌         |          ✅          |
|                             | **Algorithme de consolidation (3 matchs min. garanti)** |       ❌        |     ❌     |        ❌         |   **✅ Exclusif**    |
| **Développement & Tech**    | Grilles d'évaluation par calibre (Technique/Mental)     |       ❌        |     ❌     |    ⚠️ Basique     |     **✅ Natif**     |
|                             | **Architecture Multi-Sports progressive (`sport_id`)**  |   ⚠️ En silos   | ⚠️ Basique | ⚠️ Soccer-centric |   **✅ Abstraite**   |
|                             | **Mode Hors-Ligne Terrain (PWA)**                       |       ❌        | ⚠️ Partiel |    ⚠️ Partiel     | **✅ Offline-First** |

---

## 3. Facteurs Clés de Différenciation (USP)

1. **Flexibilité Opérationnelle :** Capacité de mélanger des catégories (ex: U9-U10) pour la saison régulière et de créer des équipes par niveau uniquement pour les tournois sans casser l'historique du sportif.
2. **Transparence Financière par Événement :** Isolation complète des coûts de chaque tournoi (inscriptions, autobus, hôtel, paie d'arbitres) face aux revenus des sous-inscriptions.
3. **Expérience Parent Maximisée :** Le parent voit l'horaire réel de son enfant, peu importe dans quel sous-groupe éphémère il est affecté cette semaine-là.

# ANALYSE D'AFFAIRES — PROJECT PULSE (PULSE OS)

## 1. Vision du Produit

Plateforme SaaS universelle de gestion de **sports collectifs** (focus initial sur le **Soccer** en Phase 1, préparée pour le hockey, basket-ball et autres en Phase 2).
PElle couvre le secteur récréatif, le développement et le compétitif.
Elle permet d'inscrire des enfants à des bassins d'âge, de composer des équipes éphémères (saison, tournois, pratiques) et d'offrir un calendrier dynamique aux parents.

## 2. Rôles et Perspectives

* **Admin Club / Directeur Technique :**
  * Gestion de la structure, des terrains,
  * création des événements,
  * répartition des équipes éphémères,
  * finances.
* **Responsable d'équipe & Entraîneurs :**
  * Suivi des effectifs,
  * fiches d'évaluations (Technique, Tactique, Physique, Mental),
  * logistique terrain.
* **Représentant Familial :** Inscription des enfants,
  * sélection des sous-événements (tournois à la carte),
  * suivi du calendrier familial consolidé,
  * Inscriptions (Bassin, événements, etc.)
  * RSVP.
* **Sportif (Joueur) :**
  * Si Majeur (18+ ans)
    * Inscriptions (Bassin, événements, etc.)
  * Consultation de son alignement,
  * accès à son fichier de performance et ses objectifs.
  * suivi de son calendrier
  
## 3. Modèle Opérationnel Clé

* **Bassins d'âges (Pools) :** Groupe global d'inscrits (ex: U9-U10, Senior).
* **Groupes d'entraînement & Équipes de saison :** Structures de club persistantes à l'année (indépendantes de tout événement) — voir `ADR-008`.
* **Alignements :** Sous-groupes ponctuels formés pour un événement précis, à partir des sportifs inscrits.
* **Sous-Inscriptions :** Inscription principale à la saison + options à la carte (Tournois #1, #3, #4).
* **Organisations :** Un club peut être rattaché à une association régionale ou une fédération (imbrication arbitraire) — voir `ADR-008 §6`.

### 3.1 Flux d'un Événement (Tournoi)

1. **Un club organise** un événement — modalités (format, classement) et éligibilité (plages âge/genre, potentiellement plusieurs par événement, ex: U4-U8 et U13-M et plus).
2. **Des clubs s'inscrivent** à l'événement en y engageant une ou plusieurs équipes par catégorie éligible (l'organisateur ou d'autres clubs).
3. **Les Représentants inscrivent leurs enfants** à l'événement — une sous-inscription distincte de l'engagement du club.
4. **Le Directeur Technique construit les alignements** à partir des inscrits, un par équipe engagée, et assigne les entraîneurs pour l'occasion.
5. **Le club organisateur planifie** les matchs, terrains et arbitres (module `scheduling`).

## [TODO AMÉLIORATION AFFAIRES]

* [ ] Réfléchir à la gestion des bordereaux d'affiliation avec la fédération.
* [ ] Définir les critères de calcul automatique pour les remboursements en cas d'annulation de tournoi.

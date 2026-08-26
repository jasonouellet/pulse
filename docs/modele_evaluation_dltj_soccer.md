# Modèle d'Évaluation et d'Architecture pour la Gestion Sportive (DLTJ / LTPD)

## 1. Contexte & Philosophie d'Évaluation

Ce modèle s'appuie sur le cadre de **Développement à long terme du joueur (DLTJ / LTPD)** promu par **Soccer Canada**, **Soccer Québec** (notamment dans le cadre des Centres de Développement de Club - CDC) et inspiré des meilleures académies européennes (modèle TIPS d'Ajax, matrice 4 coins UEFA/FA).

### Principes Directeurs

- **Développement continu (No Cut Policy) :** Remplacement des essais uniques (_tryouts_) par une observation régulière en entraînement et en match.
- **Prévention du biais de maturité biologique :** Séparation du potentiel athlétique/technique de la simple taille ou force brute. Le pilier physique mesure le savoir-faire moteur et l'agilité, et non l'impact physique brut.
- **Équilibration dynamique :** Utilisation des scores pondérés pour créer des groupes et équipes homogènes (répartition en serpentin), garantissant un temps de jeu et un défi adaptés à chaque enfant.
- **Double niveau d'observation (Terrain vs Bilan) :** Un filtre `is_quick_eval` permet de basculer entre une observation rapide sur le terrain (10 critères clés) et une évaluation bilan complète (28 critères).

---

## 2. Grille des 28 Critères d'Évaluation (5 Piliers)

### 1. Technique & Tactique Individuelle (25% en U10-U12)

- **Maîtrise & Conduite de balle** `[Terrain]` : Capacité à avancer avec le ballon collé au pied à différentes vitesses.
- **Contrôle orienté** `[Terrain]` : Prise de balle qui élimine un adversaire ou prépare l'action suivante.
- **Jeu court (Passes)** : Précision, dosage et timing des passes au sol (intérieur du pied).
- **Jeu long & Frappes** : Capacité à lever le ballon, renverser le jeu ou tirer avec puissance et précision.
- **Duels offensifs (1v1)** `[Terrain]` : Élimination par la feinte, le changement de rythme ou la protection de balle.
- **Duels défensifs (1v1)** : Cadrage, freinage de l'adversaire et tacle/interception propre.
- **Prise d'information (Scanning)** `[Terrain]` : Fréquence des coups d'œil pour observer l'environnement avant réception.
- **Prise de décision** `[Terrain]` : Choix de la meilleure option selon la situation (ex: repasser par l'arrière vs dribbler).

### 2. Tactique Collective & Vision du Jeu (15% en U10-U12)

- **Occupation du terrain** : Utilisation de la largeur et de la profondeur selon la phase de jeu.
- **Offre de solution (Soutien)** `[Terrain]` : Se déplacer continuellement dans un espace libre pour aider le porteur.
- **Compréhension des rôles** : Respect des principes de sa position sans se télescopier avec les coéquipiers.
- **Transition Attaque -> Défense** : Vitesse de réaction à la perte du ballon (gagner du temps, cadrer).
- **Transition Défense -> Attaque** : Vitesse de projection ou de conservation dès la récupération.

### 3. Physique & Motricité (20% en U10-U12)

- **Agilité & Équilibre** `[Terrain]` : Qualité des appuis et changements de direction brusques sans perte d'équilibre.
- **Coordination motrice** : Fluidité globale du mouvement et dissociation haut/bas du corps.
- **Vitesse / Explosivité** `[Terrain]` : Réactivité et accélération sur les 3 à 5 premiers mètres.
- **Vitesse maximale** : Vitesse de pointe sur une longue course de repli ou de contre-attaque.
- **Endurance active** : Capacité à maintenir un rythme élevé tout au long du match/entraînement.
- **Aisance du pied faible** : Niveau de coordination et d'efficacité avec le pied non dominant.

### 4. Mental & Cognitif (20% en U10-U12)

- **Concentration & Écoute** `[Terrain]` : Rétention des consignes et capacité à rester focalisé sur l'exercice.
- **Résilience (Gestion de l'erreur)** `[Terrain]` : Réaction après avoir perdu la balle (effort immédiat vs frustration).
- **Intensité & Combativité** `[Terrain]` : Volonté de s'investir pleinement dans chaque duel et chaque atelier.
- **Autonomie & Initiative** : Capacité à essayer des choses créatives sans attendre l'ordre du coach.
- **Gestion de la pression** : Niveau de jeu préservé lors de matchs à enjeu ou sous forte opposition.

### 5. Socio-Émotionnel (20% en U10-U12)

- **Communication positive** `[Terrain]` : Encouragements explicites, absence de reproches envers les partenaires.
- **Coopération & Esprit d'équipe** : Plaisir de faire passer le collectif avant sa statistique personnelle.
- **Respect & Étiquette** `[Terrain]` : Attitude envers l'arbitre, l'adversaire, le matériel et les éducateurs.
- **Assiduité & Ponctualité** : Présence, rigueur dans la préparation et respect des règles du groupe.

---

## 3. Matrice de Pondération par Tranche d'Âge

| Pilier                          | U7 à U9 (Fondamentaux) | U10 à U12 (Apprendre à s'entraîner) | U13+ (S'entraîner à compétitionner) |
| :------------------------------ | :--------------------: | :---------------------------------: | :---------------------------------: |
| **Technique & Tactique Indiv.** |          30 %          |                25 %                 |                25 %                 |
| **Tactique Collective**         |          10 %          |                15 %                 |                20 %                 |
| **Physique (Motricité)**        |          20 %          |                20 %                 |                25 %                 |
| **Mental / Cognitif**           |          20 %          |                20 %                 |                20 %                 |
| **Socio-Émotionnel**            |          20 %          |                20 %                 |                10 %                 |

---

## 4. Architecture de Données Générique et Évolutive (Multisport)

```text
+------------------+       +-------------------------+       +--------------------+
|      Sport       |1    * |   Evaluation_Framework  |1    * |       Pillar       |
|------------------|-------|-------------------------|-------|--------------------|
| id               |       | id                      |       | id                 |
| name             |       | sport_id                |       | framework_id       |
+------------------+       | name                    |       | name               |
                           | age_min / age_max       |       | weighting_percent  |
                           +-------------------------+       +--------------------+
                                                                        | 1
                                                                        | *
                                                             +--------------------+
                                                             |     Criterion      |
                                                             |--------------------|
                                                             | id                 |
                                                             | pillar_id          |
                                                             | name               |
                                                             | is_quick_eval (bool)|
                                                             | scale_min / max    |
                                                             +--------------------+
```

### 1. Tables de Configuration (Métadonnées Moteur)

- **`Sport`** : `id`, `name` (ex: Soccer, Hockey, Basketball).
- **`Evaluation_Framework`** : `id`, `sport_id`, `name`, `age_min`, `age_max` (ex: DLTJ CDC U10-U12).
- **`Pillar`** : `id`, `framework_id`, `name`, `weighting_percent` (ex: Technique Indiv. 25%).
- **`Criterion`** : `id`, `pillar_id`, `name`, `is_quick_eval`, `scale_min`, `scale_max` (ex: Scanning, boolean `true`, Échelle 1-5).

### 2. Tables Opérationnelles & Suivi

- **`Athlete`**, **`Season_Group`**, **`Roster`**, **`Coach`**.
- **`Observation_Session`** : `id`, `group_id`, `coach_id`, `date`, `context` (Entraînement, Match, Bilan).
- **`Evaluation_Score`** : `id`, `session_id`, `athlete_id`, `criterion_id`, `score_value`, `notes`.

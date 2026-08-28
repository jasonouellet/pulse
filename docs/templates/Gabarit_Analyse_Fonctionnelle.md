# Gabarit d'Analyse Fonctionnelle (Dossier de Cadrage & Exigences)

Ce document constitue un gabarit d'analyse fonctionnelle structuré, aligné sur les meilleures pratiques de l'architecture d'affaires et de l'analyse d'affaires:

* **BIZBOK®** (Business Architecture Body of Knowledge) : Alignement stratégique, capacités et chaîne de valeur.
* **BABOK® v3** (Business Analysis Body of Knowledge) : Ingénierie des exigences, cycle de vie et traçabilité.
* **Strategyzer** : Business Model Canvas & Value Proposition Canvas.
* **Value Stream Mapping (VSM)** : Cartographie et optimisation du flux de valeur.

---

## 1. Contextualisation Stratégique & Alignement (BIZBOK® & Strategyzer)

### 1.1 Alignement Stratégique & Modèle d'Affaires

* **Objectif stratégique :** *(Décrire l'alignement avec la vision et les objectifs stratégiques de l'organisation)*
* **Business Model Canvas (Synthèse) :**
  * **Proposition de valeur :** *(Ce que la solution apporte aux clients/utilisateurs)*
  * **Segments de clientèle / Intéressés :** *(Qui bénéficie directement de cette évolution)*
  * **Partenaires & Ressources clés :** *(Équipes, dépendances système ou prestataires requis)*

### 1.2 Alignement de la Valeur (Value Proposition Canvas)

* **Profil Client / Utilisateur (*Customer Profile*) :**
  * **Tâches (*Customer Jobs*) :** *(Ce que l'utilisateur essaie d'accomplir au quotidien)*
  * **Problèmes (*Pains*) :** *(Frustrations, obstacles, risques actuels)*
  * **Gains attendus (*Gains*) :** *(Bénéfices souhaités ou mesurables)*
* **Carte de Valeur (*Value Map*) :**
  * **Produits & Services :** *(La solution proposée)*
  * **Analgésiques (*Pain Relievers*) :** *(Comment la solution supprime ou réduit les problèmes)*
  * **Générateurs de gains (*Gain Creators*) :** *(Comment la solution apporte de la valeur mesurable)*

### 1.3 Chaîne de Valeur & Processus (Value Stream Mapping - VSM)

* **Chaîne de valeur globale (BIZBOK® Value Stream) :** *(Ex. De la prise de besoin client jusqu'à la livraison du service)*
* **Cartographie VSM du processus cible :**
  * **État Actuel (*Current State*) vs État Cible (*Future State*)**
  * **Indicateurs clés VSM :**
    * Temps de cycle (*Cycle Time - CT*)
    * Temps à valeur ajoutée vs Temps de non-valeur ajoutée (*VA / NVA*)
    * Goulots d'étranglement identifiés et opportunités d'automatisation.

---

## 2. Analyse du Périmètre & Délimitation (BABOK® - Scope & Elicitation)

### 2.1 Analyse des Parties Prenantes (*Stakeholder Analysis*)

| Rôle / Titre               | Entité / Département | Niveau d'impact | Rôle dans le projet (RACI)      |
| :------------------------- | :------------------- | :-------------- | :------------------------------ |
| *Ex: Gestionnaire de paie* | *RH*                 | *Élevé*         | *A (Approbateur / Accountable)* |
| *Ex: Utilisateur final*    | *Ventes*             | *Moyen*         | *C (Consulté / Consulted)*      |

### 2.2 Délimitation du Périmètre (*In-Scope / Out-of-Scope*)

* **Dans le périmètre (*In-Scope*) :** *(Liste explicite des fonctionnalités, modules et systèmes inclus)*
* **Hors périmètre (*Out-of-Scope*) :** *(Liste explicite des éléments délibérément exclus pour cette phase)*

---

## 3. Spécification des Exigences (BABOK® Requirements Standard)

### 3.1 Exigences d'Affaires (*Business Requirements*)

* **R-BUS-01 :** *(Ex. Réduire le temps de traitement des demandes de 30 % d'ici la fin du T4)*

### 3.2 Exigences Utilisateurs (*User Requirements & User Stories*)

* **Format recommandé :** *En tant que [Rôle], je veux [Fonctionnalité], afin de [Bénéfice/Valeur].*

* **Exemple de critère d'acceptation (Format Gherkin / BDD) :**

  ```gherkin
  Étant donné que [Mise en situation initiale]
  Lorsque [Action effectuée par l'utilisateur]
  Alors [Résultat attendu du système]
  ```

### 3.3 Exigences Fonctionnelles (*Functional Requirements*)

| ID Exigence | Module / Composant    | Description de la règle fonctionnelle                                | Priorité (MoSCoW) |
| :---------- | :-------------------- | :------------------------------------------------------------------- | :---------------- |
| **RF-01**   | *Gestion des comptes* | *Le système doit valider le format du courriel avant la soumission.* | *Must Have*       |
| **RF-02**   | *Rapports & Exports*  | *Le système doit permettre l'exportation des données au format CSV.* | *Should Have*     |

### 3.4 Exigences Non-Fonctionnelles (*Non-Functional Requirements*)

* **Performance & Scalabilité :** *(Ex. Temps de réponse < 2 secondes pour 95% des requêtes)*
* **Sécurité & Conformité :** *(Ex. Chiffrement des données sensibles au repos et en transit)*
* **Ergonomie & Accessibilité (UX/UI) :** *(Ex. Respect des normes WCAG 2.1 AA)*
* **Disponibilité & Fiabilité :** *(Ex. Taux de disponibilité cible de 99,9 %)*

---

## 4. Modélisation Fonctionnelle & Données

### 4.1 Diagrammes de Processus (BPMN)

> (Insérer ou référencer les schémas de processus de niveau 2 ou 3)

### 4.2 Modèle de Données & Dictionnaire de Données

* **Modèle Conceptuel de Données (MCD) / Modèle d'Entités**

* **Dictionnaire de Données :**

| Champ             | Type de donnée | Obligatoire (O/N) | Valeurs autorisées / Règles de gestion    |
| :---------------- | :------------- | :---------------- | :---------------------------------------- |
| `Statut_Commande` | *Texte / Enum* | *Oui*             | *[Brouillon, Validée, Expédiée, Annulée]* |

---

## 5. Matrice de Traçabilité & Validation (BABOK® Life Cycle Management)

| Besoin d'Affaires (BIZBOK / Strategyzer) | Exigence Fonctionnelle (BABOK) | Cas d'utilisation / User Story | Cas de Test (Recette) |
| :--------------------------------------- | :----------------------------- | :----------------------------- | :-------------------- |
| *R-BUS-01*                               | *RF-01*                        | *US-01*                        | *TC-01*               |

---

## 6. Approbation & Validation (*Sign-Off*)

| Nom  | Rôle                    | Statut (Approuvé / En attente) | Date |
| :--- | :---------------------- | :----------------------------- | :--- |
|      | *Sponsor Projet*        |                                |      |
|      | *Architecte d'Affaires* |                                |      |
|      | *Analyste Fonctionnel*  |                                |      |

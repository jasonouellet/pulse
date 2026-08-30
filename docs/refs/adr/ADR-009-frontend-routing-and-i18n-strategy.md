# [ADR-009] Stratégie d'internationalisation (i18n) et de routage pour le frontend

* Status: accepted
* Deciders: PULSE Core Engineering Team
* Date: 2026-08-30

## Contexte et énoncé du problème

L'application frontend de Project PULSE doit servir une base d'utilisateurs bilingue (Français/Anglais) comprenant des administrateurs, entraîneurs et parents. L'interface nécessitait une architecture fluide capable de gérer le changement dynamique de langue, la mise en page imbriquée (*nested layouts*) axée sur les rôles utilisateur, tout en garantissant un typage strict TypeScript et de bonnes performances de rechargement à chaud avec Vite 6 et Node.js v22.

Comment structurer de manière pérenne l'internationalisation et la navigation dans l'application tout en évitant la duplication de code et les erreurs de typage à l'exécution ?

## Facteurs décisionnels

* **Type Safety :** Nécessité d'avoir un typage strict TypeScript sur les clés de traduction pour éviter les clés manquantes ou obsolètes lors de la refactorisation.
* **Architecture de Layouts :** Séparation claire entre les wrappers applicatifs (Header, Navigation, Theme) et les composants de pages métiers.
* **Bilinguisme dynamique :** Prise en charge native du Français (langue par défaut) et de l'Anglais sans rechargement de page.
* **Maintenabilité :** Alignement strict entre les modèles de données (ex: rôles `UserRole`) et les dictionnaires i18n.
* **Compatibilité Tooling :** Intégration optimale avec React Router v7 / React 19 et Vite 6.

## Options envisagées

* **Option 1 :** React Router avec le composant `<Outlet/>` + `i18next` avec typage strict via `CustomTypeOptions`
* **Option 2 :** Passage explicite de `children` dans `AppLayout` + `react-intl` (FormatJS)
* **Option 3 :** Routage personnalisé basé sur des états locaux + fichier de dictionnaire JS maison

## Décision retenue

Option choisie : **Option 1**, car elle offre la meilleure modularité pour les layouts imbriqués via le pattern déclaratif de React Router, couplée à l'écosystème robuste et éprouvé d'`i18next`.

### Conséquences positives

* **Pattern Layout Propre :** `AppLayout` agit comme route parente englobant l'en-tête, le sélecteur de thème, le composant `<LanguageSwitcher/>` et le menu utilisateur, déléguant l'affichage des sous-pages au composant `<Outlet/>`.
* **Synchronisation stricte des Types :** Les clés de rôle `UserRole` en majuscules (`"CLUB_ADMIN"`, `"TECHNICAL_DIRECTOR"`, etc.) sont calquées à 1:1 dans le dictionnaire JSON (`common.json`), garantissant un typage sécurisé sans transformation manuelle de chaînes.
* **Détection et Persistance :** Intégration de `i18next-browser-languagedetector` permettant la persistance de la langue choisie dans le stockage local du navigateur.

### Conséquences négatives

* **Légère verbosité initiale :** Obligation de maintenir des fichiers JSON jumeaux (`fr/common.json`, `en/common.json`) sous le même schéma structuré.

## Arguments pour et contre les options

### Option 1 : React Router `<Outlet/>` + `i18next` (Option retenue)

* **Pour :** Découplage complet entre les routes et la mise en page globale.
* **Pour :** Support natif du typage automatique des clés via augmentation du module `i18next`.
* **Pour :** Intégration simple du composant de bascule de langue `<LanguageSwitcher/>`.
* **Contre :** Nécessite une rigueur sur la casse des clés JSON pour correspondre aux enums/types du projet.

### Option 2 : Passage de `children` + `react-intl`

* **Pour :** Très populaire pour les applications à grande échelle.
* **Contre :** Nécessite d'envelopper chaque page manuellement dans le layout si des routes secondaires ou des sous-panneaux d'administration sont introduits.
* **Contre :** Syntaxe de composants (`<FormattedMessage/>`) parfois plus lourde à lire que le hook `t()`.

### Option 3 : Solution maison

* **Pour :** Aucun package externe requis.
* **Contre :** Réinvention de la roue pour la gestion du pluriel, la détection de la langue du navigateur et la persistance.
* **Contre :** Risque élevé de bugs de régression lors de la montée en charge du projet.

## Validation de la décision

La décision est considérée comme validée car :

1. L'application bascule instantanément entre le français et l'anglais sans rechargement de page.
2. Les éléments de navigation dynamiques (`nav`), les informations de session utilisateur (`user`) et l'intégralité des rôles (`roles`) s'affichent correctement dans la langue sélectionnée.
3. La compilation TypeScript (`npx tsc --noEmit`) s'exécute sans aucun avertissement ni erreur de typage.

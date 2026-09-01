# [ADR 001-B]: Adoption des Primitives Headless (Radix UI) et Tailwind CSS pour le Design System

## Statut

Accepté (Précise et amende l'ADR 001)

## Contexte

L'interface du **Project PULSE** doit desservir plusieurs types d'utilisateurs (parents sur mobile, entraîneurs sur le terrain, administrateurs de club sur PC).

Nous avions besoin d'une stratégie UI qui garantisse :

1. Une accessibilité sans faille (norme WCAG 2.1 AA natif, lecteurs d'écran, navigation au clavier).
2. Une flexibilité totale sur le design (liberté de créer une expérience fluide pour le sport sans subir le style pré-imposé d'une bibliothèque tierce).
3. Une empreinte JavaScript minimale pour une exécution ultra-rapide sur le réseau mobile du terrain.

Nous devions choisir entre les bibliothèques UI monolithiques classiques (ex: Material UI, Ant Design, Bootstrap) et l'approche moderne par primitives "Headless" (Radix UI + Tailwind CSS).

## Décision

Nous adoptons le couple **Radix UI (Primitives Headless)** et **Tailwind CSS**.

### 1. Pourquoi Radix UI (Headless Primitives) ?

* **Découplage Logique / Style :** Radix fournit la mécanique complexe (gestion du focus, états ouvert/fermé, événements clavier, attributs WAI-ARIA, fermetures sur clic extérieur) sans appliquer aucun style CSS.
* **Accessibilité Natifs (a11y) :** Tous les composants respectent la norme WCAG 2.1 AA sans effort de développement supplémentaire.
* **Architecture Modulaire :** Chaque composant est un paquet indépendant (ex: `@radix-ui/react-dialog`), ce qui évite d'importer une bibliothèque lourde inopérante.

### 2. Pourquoi Tailwind CSS ?

* **Utility-First :** Permet de composer le design directement sur les primitives Radix UI sans écrire de CSS surchargé.
* **Compilation Optimisée (Tree-shaking) :** Seules les classes CSS réellement utilisées sont incluses dans le build final (taille de bundle minimale).
* **Prise en charge du Dark/Light Mode :** Gestion native via des variables CSS et des classes utilitaires (`dark:`).

## Conséquences

### Impacts Positifs

* **Standard du Marché :** Alignement complet sur les meilleures pratiques contemporaines de l'écosystème React (fondation de Shadcn UI).
* **Contrôle Visuel Total :** Aucune guerre de surcharge CSS (`!important` ou styles inline) pour personnaliser les composants.
* **Expérience Terrain Optimisée :** Interface légère et rapide à charger sur les réseaux cellulaires.

### Compromis / Contraintes

* Demande de composer soi-même la couche visuelle des éléments de base (boutons, modales) lors de la création du Design System local (`src/components/ui/`), contrairement à une solution "clef en main" comme Material UI. Ce surcoût initial est compensé par la réutilisabilité et la maintenance à long terme.

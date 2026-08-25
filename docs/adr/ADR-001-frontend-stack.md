# ADR 001 : Choix de la Pile Technologique UI/UX & Support Multi-Plateforme

## Statut

Accepté

## Contexte

L'application doit supporter plusieurs profils (Parents, Sportifs, Administrateurs) sur mobiles, tablettes et PC, offrir une accessibilité irréprochable (WCAG 2.1 AA), supporter les thèmes Clair/Sombre et offrir une expérience utilisateur optimale lors d'une utilisation sur le terrain.

## Décisions

1. **Framework UI :** React.js (TypeScript) avec approche Mobile-First et Responsive Design (Breakpoints Tailwind CSS : `sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`).
2. **Système de Design & Composants :** Radix UI / Shadcn UI (Primitives headless 100% accessibles WAI-ARIA natif) couplé à Tailwind CSS.
3. **Gestion des Thèmes (Dark/Light) :** CSS Variables avec détection dynamique de la préférence système (`prefers-color-scheme`) et basculement manuel (stockage local).
4. **Accessibilité (a11y) :** Respect strict de la norme WCAG 2.1 AA (contraste minimal de 4.5:1, navigation au clavier complète, zones de touche mobiles de minimum 44x44px, attributs ARIA systématiques).
5. **Tactile & PWA :** Support PWA (Progressive Web App) pour une utilisation offline partagée sur le terrain et gestes tactiles optimisés (Swipe/Touch target).

## Conséquences

- Développement accéléré de composants accessibles grâce aux primitives Radix UI.
- Performance élevée et empreinte mémoire faible adaptée au réseau mobile du terrain.

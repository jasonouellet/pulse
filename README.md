# PULSE

Plateforme SaaS universelle de gestion de sports collectifs (multi-niveaux, bassins d'âges, rosters éphémères, tournois et calendrier dynamique).

## 🚀 Démarrage Rapide (Docker / Podman)

1. Copier les variables d'environnement :

   ```bash
   cp .env.example .env
   ```

2. Lancer l'infrastructure complète:

   ```bash
   docker compose up -d
   # ou avec Podman :
   # podman-compose up -d
   ```

3. Vérifier la santé du backend :

   ```bash
   curl http://localhost:8080/healthz
   ```

## 📚 Documentation

- `docs/SYSTEM_CONTEXT.md` : Analyse d'affaires et cas d'usage.
- `docs/ARCHITECTURE.md` : Choix techniques et découpage des modules.
- `docs/ROADMAP_TODO.md` : Suivi de l'avancement.
- `docs/adr/` : Registre des décisions d'architecture (ADRs 001 à 005).
- `LICENSE` & `CLA.md` : Licence BSL Non-Commerciale et accord de contribution.

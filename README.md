# PULSE

Universal SaaS platform for managing team sports (multi-level, age pools, ephemeral rosters, tournaments, and dynamic scheduling).

## 🚀 Quick Start (Docker / Podman)

1. Copy the environment variables:

   ```bash
   cp .env.example .env
   ```

2. Start the full infrastructure:

   ```bash
   docker compose up -d
   # or with Podman:
   # podman-compose up -d
   ```

3. Check the backend health:

   ```bash
   curl http://localhost:8080/healthz
   ```

## 📚 Documentation

- `docs/SYSTEM_CONTEXT.md`: Business analysis and use cases.
- `docs/ARCHITECTURE.md`: Technical choices and module breakdown.
- `docs/ROADMAP_TODO.md`: Progress tracking.
- `docs/adr/`: Architecture decision records (ADR 001 to 005).
- `LICENSE` & `CLA.md`: Non-Commercial BSL license and contributor agreement.

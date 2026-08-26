FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm

# 1. Activation du téléchargement automatique de la chaîne d'outils Go si une dépendance exige une version plus récente
ENV GOTOOLCHAIN=auto

# 2. Installation des dépendances système & Node.js avec nettoyage propre du cache apt
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get update \
    && apt-get install -y --no-install-recommends nodejs postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# 3. Installation des outils Go
# Note : Pour figer une version compatible Go 1.25 sans dépendre de GOTOOLCHAIN=auto, 
# vous pouvez remplacer @latest par @v1.28.0
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/air-verse/air@latest

# 4. Installation de golang-migrate (extraction directe vers /usr/local/bin)
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar -xz -C /usr/local/bin migrate

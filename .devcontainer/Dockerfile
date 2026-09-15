FROM mcr.microsoft.com/devcontainers/go:1.26-bookworm

# 1. Activation du téléchargement automatique de la chaîne d'outils Go si une dépendance exige une version plus récente
ENV GOTOOLCHAIN=auto

# 2. Installation des dépendances système & Node.js avec nettoyage propre du cache apt
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl -fsSL https://deb.nodesource.com/setup_24.x | bash - \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
    postgresql-client \
    python3-pip \
    python3-venv \
    pipx \
    git-lfs \
    && git lfs install \
    && rm -rf /var/lib/apt/lists/*

# Installations système pip et binaries
RUN pip install --no-cache-dir --break-system-packages pre-commit sqlfluff

# Fixer la propriété du dossier GOPATH dès le départ pour l'utilisateur vscode
RUN mkdir -p /go && chown -R vscode:vscode /go

# Passer à l'utilisateur vscode pour installer les outils Go dans le bon profil
USER vscode

# 3. Installation des outils Go
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/air-verse/air@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/dependabot/cli/cmd/dependabot@latest

# 5. Installation de golang-migrate (extraction directe vers /usr/local/bin)
# Version alignée avec github.com/golang-migrate/migrate/v4 dans go.mod
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar -xz -C /usr/local/bin migrate

# Assurer l'accès final aux répertoires Go
RUN chown -R vscode:vscode /go
